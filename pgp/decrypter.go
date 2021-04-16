package pgp

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type Decrypter struct {
	originalReader io.Reader
	armored        bool
	recipients     *KeyRing
	signers        *KeyRing
	symKey         *SymKey
	maxTempMemory  int

	armorBlock     *armor.Block
	reader         io.Reader
	msg            *openpgp.MessageDetails
	tempKey        *SymKey
	tempFile       string
	tempFileReader io.ReadCloser
}

func CreateDecrypter(reader io.Reader) *Decrypter {
	return &Decrypter{originalReader: reader, maxTempMemory: 1024 * 1024}
}

func (me *Decrypter) MaxTempMemory(maxTempMemory int) *Decrypter {
	me.maxTempMemory = maxTempMemory
	return me
}

func (me *Decrypter) Armored(armored bool) *Decrypter {
	me.armored = armored
	return me
}

func (me *Decrypter) Decrypt(recipients *KeyRing) *Decrypter {
	me.recipients = recipients
	return me
}

func (me *Decrypter) Signers(signers *KeyRing) *Decrypter {
	me.signers = nil
	if signers != nil {
		me.signers = signers.PubOnly()
	}
	return me
}

func (me *Decrypter) Symmetric(key *SymKey) *Decrypter {
	me.symKey = key
	return me
}

func (me *Decrypter) Start() io.Reader {
	util.Assert(me.originalReader == nil, "Reader is required")
	util.Assert(me.symKey == nil && me.recipients == nil && me.signers == nil && !me.armored, "Nothing to do")
	util.Assert(me.symKey != nil && me.recipients != nil, "Symmetric decryption can not have recipients")
	util.Assert(me.symKey != nil && me.signers != nil, "Symmetric decryption can not have signers")

	me.reader = me.originalReader
	me.preapreArmored()

	if me.symKey == nil && me.recipients == nil && me.signers == nil {
		log.Printf("Decrypter, armor parsing only")
		return &decryptor2Reader{decrypter: me}
	}
	if me.symKey != nil {
		return me.openSymDecrypt()
	}

	me.unsafeDecrypt()
	util.Assert(me.recipients != nil && !me.msg.IsEncrypted, "Decrypt, it is not encrypted")
	util.Assert(me.recipients != nil && me.msg.IsSymmetricallyEncrypted, "Decrypt, it is symmetrically encrypted")

	if me.signers == nil {
		return me.reader
	}

	me.decryptToTemp()
	me.checkSign()
	return me.openTemp()
}

func (me *Decrypter) openTemp() io.Reader {
	if me.tempFile == "" {
		return me.reader
	}
	ret, err := os.Open(me.tempFile)
	util.Check(err)
	me.tempFileReader = ret
	me.reader = CreateDecrypter(ret).Symmetric(me.tempKey).Start()
	return me.reader
}

func (me *Decrypter) checkSign() {
	if me.signers == nil {
		return
	}
	util.Assert(!me.msg.IsSigned, "Decrypt, msg is not signed")
	util.Assert(me.msg.Signature == nil, "Decrypt, unknown signer: %X", me.msg.SignedByKeyId)
	sigKP := keyFromEntity(me.msg.SignedBy.Entity)
	pubKey := sigKP.ExportPubArmored()
	for _, v := range me.signers.kps {
		if v.ExportPubArmored() == pubKey {
			return
		}
	}
	log.Panicf("Decrypt, signer is not a writer: %X %s", sigKP.Id(), sigKP.UserName())
}

func (me *Decrypter) openSymDecrypt() io.Reader {
	msg, err := openpgp.ReadMessage(me.reader, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return me.symKey.key, nil
	}, Config)
	util.Check(err)
	me.msg = msg
	log.Printf("Decrypt, symmetric with key size: %d", me.symKey.Size())
	me.reader = me.msg.UnverifiedBody
	return &decryptor2Reader{decrypter: me}
}

func (me *Decrypter) preapreArmored() {
	if !me.armored {
		return
	}
	log.Printf("Decrypter, Prepering armor")
	block, err := armor.Decode(me.reader)
	util.Check(err)
	me.armorBlock = block
	me.reader = block.Body
}

func (me *Decrypter) unsafeDecrypt() {
	keys := KeyRingCreate().toPgpEntityList()
	if me.signers != nil {
		keys = append(keys, me.signers.toPgpEntityList()...)
	}
	if me.recipients != nil {
		keys = append(keys, me.recipients.toPgpEntityList()...)
	}
	msg, err := openpgp.ReadMessage(me.reader, keys, nil, nil)
	util.Check(err)
	me.msg = msg
	me.reader = msg.UnverifiedBody
}

func (me *Decrypter) decryptToTempMemory() []byte {
	buf := make([]byte, me.maxTempMemory+1)
	read, err := io.ReadAtLeast(me.reader, buf, me.maxTempMemory+1)
	if err == io.ErrUnexpectedEOF {
		return buf[0:read]
	}
	util.Check(err)
	if read != me.maxTempMemory+1 {
		log.Panicf("Wrong, expect: %d, but was: %d", me.maxTempMemory+1, read)
	}
	return buf
}

func (me *Decrypter) decryptToTemp() {
	buf := me.decryptToTempMemory()
	if len(buf) <= me.maxTempMemory {
		log.Printf("Decrypt to memory, total: %d", len(buf))
		me.reader = bytes.NewBuffer(buf)
		return
	}

	f, err := ioutil.TempFile(os.TempDir(), "vaultz-decrypt-*.tmp")
	util.Check(err)
	defer f.Close()
	me.tempKey = SymKeyGenerate()
	encrypter := CreateEncrypter(f).Symmetric(me.tempKey)
	defer encrypter.Close()
	w := encrypter.Start()
	w.Write(buf)
	total, err := io.Copy(w, me.reader)
	util.Check(err)
	log.Printf("Decrypt to %s, total: %d", f.Name(), total+int64(len(buf)))
	me.tempFile = f.Name()
}

type decryptor2Reader struct {
	decrypter *Decrypter
}

func (me *decryptor2Reader) Read(p []byte) (n int, err error) {
	return me.decrypter.reader.Read(p)
}

func (me *Decrypter) Close() error {
	if me.tempFileReader != nil {
		log.Printf("Decrypter closing, closing temp file: %s", me.tempFile)
		me.tempFileReader.Close()
	}
	if me.tempFile != "" {
		log.Printf("Decrypter closing, deleting file: %s", me.tempFile)
		os.Remove(me.tempFile)
	}
	return nil
}
