package pgp

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type Decrypter2 struct {
	originalReader io.Reader
	armored        bool
	recipients     *KeyRing
	writers        *KeyRing
	symKey         *SymKey

	armorBlock     *armor.Block
	reader         io.Reader
	msg            *openpgp.MessageDetails
	tempKey        *SymKey
	tempFile       string
	tempFileReader io.ReadCloser
}

type Decryptor2Reader struct {
	decrypter *Decrypter2
}

func (me *Decryptor2Reader) Read(p []byte) (n int, err error) {
	return me.decrypter.reader.Read(p)
}

func (me *Decrypter2) Close() error {
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

func CreateDecrypter(reader io.Reader) *Decrypter2 {
	return &Decrypter2{originalReader: reader}
}

func (me *Decrypter2) Armor(armored bool) *Decrypter2 {
	me.armored = armored
	return me
}

func (me *Decrypter2) Decrypt(recipients *KeyRing) *Decrypter2 {
	me.recipients = recipients
	return me
}

func (me *Decrypter2) Verify(writers *KeyRing) *Decrypter2 {
	me.writers = writers
	return me
}

func (me *Decrypter2) Symmetric(key *SymKey) *Decrypter2 {
	me.symKey = key
	return me
}

func (me *Decrypter2) Start() io.Reader {
	me.check(me.originalReader == nil, "Reader is required")
	me.check(me.symKey == nil && me.recipients == nil && me.writers == nil && !me.armored, "Nothing to do")
	me.check(me.symKey != nil && me.recipients != nil, "Symmetric decryption can not have recipients")
	me.check(me.symKey != nil && me.writers != nil, "Symmetric decryption can not have writers")

	me.reader = me.originalReader
	me.preapreArmored()

	if me.symKey == nil && me.recipients == nil && me.writers == nil {
		return me.openArmored()
	}
	if me.symKey != nil {
		return me.openSymDecrypt()
	}

	me.unsafeDecrypt()
	me.check(me.recipients != nil && !me.msg.IsEncrypted, "Decrypt, it is not encrypted")
	me.check(me.recipients != nil && me.msg.IsSymmetricallyEncrypted, "Decrypt, it is symmetrically encrypted")

	me.decryptToTemp()
	me.checkSign()

	return me.openTempFile()
}

func (me *Decrypter2) openTempFile() io.Reader {
	ret, err := os.Open(me.tempFile)
	util.Check(err)
	me.tempFileReader = ret
	me.reader = CreateDecrypter(ret).Symmetric(me.tempKey).Start()
	return me.reader
}

func (me *Decrypter2) checkSign() {
	if me.writers == nil {
		return
	}
	me.check(!me.msg.IsSigned, "Decrypt, msg is not signed")
	me.check(me.msg.Signature == nil, "Decrypt, unknown signer: %X", me.msg.SignedByKeyId)
	sigKP := keyFromEntity(me.msg.SignedBy.Entity)
	pubKey := sigKP.ExportPub()
	for _, v := range me.writers.kps {
		if v.ExportPub() == pubKey {
			return
		}
	}
	log.Panicf("Decrypt, signer is not a writer: %s", sigKP.Id())
}

func (me *Decrypter2) openArmored() io.Reader {
	log.Printf("Decrypter, armor parsing only")
	return &Decryptor2Reader{decrypter: me}
}

func (me *Decrypter2) openSymDecrypt() io.Reader {
	msg, err := openpgp.ReadMessage(me.reader, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return me.symKey.key, nil
	}, nil)
	util.Check(err)
	me.msg = msg
	log.Printf("Decrypt, symmetric with key size: %d", me.symKey.Size())
	me.reader = me.msg.UnverifiedBody
	return &Decryptor2Reader{decrypter: me}
}

func (me *Decrypter2) check(cond bool, msg string, v ...interface{}) {
	if cond {
		log.Panicf(msg, v...)
	}
}

func (me *Decrypter2) preapreArmored() {
	if !me.armored {
		return
	}
	log.Printf("Decrypter, Prepareing armor parsing")
	block, err := armor.Decode(me.originalReader)
	util.Check(err)
	me.armorBlock = block
	me.reader = block.Body
}

func (me *Decrypter2) unsafeDecrypt() {
	keys := KeyRingCreate().toPgpEntityList()
	if me.writers != nil {
		keys = append(keys, me.writers.toPgpEntityList()...)
	}
	if me.recipients != nil {
		keys = append(keys, me.recipients.toPgpEntityList()...)
	}
	msg, err := openpgp.ReadMessage(me.reader, keys, nil, nil)
	util.Check(err)
	me.msg = msg
}

func (me *Decrypter2) decryptToTemp() {
	f, err := ioutil.TempFile(os.TempDir(), "vaultz-decrypt-*.tmp")
	util.Check(err)
	defer f.Close()
	me.tempKey = SymKeyGenerate()
	encrypter := SymEncypterCreate(f, me.tempKey)
	defer encrypter.Close()
	total, err := io.Copy(encrypter.Encrypt(), me.msg.UnverifiedBody)
	util.Check(err)
	log.Printf("Decrypt to %s, total: %d", f.Name(), total)
	me.tempFile = f.Name()
}
