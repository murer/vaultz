package pgp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
)

func DecrypterCreate(plain io.Reader, writers *KeyRing, readers *KeyRing) *Decrypter {
	return &Decrypter{plain: plain, writers: writers, readers: readers, verifyOnly: false}
}

func VerifierCreate(plain io.Reader, writers *KeyRing) *Decrypter {
	return &Decrypter{plain: plain, writers: writers, verifyOnly: true}
}

type Decrypter struct {
	io.Closer
	plain   io.Reader
	writers *KeyRing
	readers *KeyRing

	msg      *openpgp.MessageDetails
	tempKey  *SymKey
	tempFile string

	verifyOnly bool
}

func (me *Decrypter) Close() error {
	if me.tempFile != "" {
		log.Printf("Decrypter closing, delete: %s", me.tempFile)
		os.Remove(me.tempFile)
	}
	return nil
}

func (me *Decrypter) UnsafeDecrypt() io.Reader {
	// ar, err := armor.Decode(me.plain)
	// util.Check(err)
	keys := me.writers.toPgpEntityList()
	if !me.verifyOnly {
		keys = append(keys, me.readers.toPgpEntityList()...)
	}
	msg, err := openpgp.ReadMessage(me.plain, keys, nil, nil)
	util.Check(err)
	me.msg = msg
	if !me.verifyOnly {
		if !me.msg.IsEncrypted {
			log.Panicf("Decrypt, it is not encrypted")
		}
		decKP := keyFromEntity(me.msg.DecryptedWith.Entity)
		log.Printf("Decrypt with: %s %s", decKP.Id(), decKP.UserName())
	}
	if !me.msg.IsSigned {
		log.Panicf("Decrypt, it is not signed")
	}
	return me.msg.UnverifiedBody
}

func (me *Decrypter) UnsafeDecryptBytes() []byte {
	r := me.UnsafeDecrypt()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	return ret
}

func (me *Decrypter) UnsafeDecryptString() string {
	return string(me.UnsafeDecryptBytes())
}

func (me *Decrypter) decryptToTemp() {
	unsafe := me.UnsafeDecrypt()
	f, err := ioutil.TempFile(os.TempDir(), "vaultz-decrypt-*.tmp")
	util.Check(err)
	defer f.Close()
	me.tempKey = SymKeyGenerate()
	encrypter := SymEncypterCreate(f, me.tempKey)
	defer encrypter.Close()
	total, err := io.Copy(encrypter.Encrypt(), unsafe)
	util.Check(err)
	log.Printf("Decrypt to %s, total: %d, verifyOnly: %t", f.Name(), total, me.verifyOnly)
	me.tempFile = f.Name()
}

func (me *Decrypter) decryptCheckSigner() {
	if me.msg.Signature == nil {
		msg := fmt.Sprintf("Decrypt unknown signer: %X", me.msg.SignedByKeyId)
		log.Panic(msg)
	}
	sigKP := keyFromEntity(me.msg.SignedBy.Entity)
	pubKey := sigKP.ExportPub()
	for _, v := range me.writers.kps {
		if v.ExportPub() == pubKey {
			return
		}
	}
	log.Panicf("Decrypt signer is not a writer: %s", sigKP.Id())
}

func (me *Decrypter) Decrypt() io.ReadCloser {
	me.decryptToTemp()
	me.decryptCheckSigner()
	signerKP := keyFromEntity(me.msg.SignedBy.Entity)
	log.Printf("Decrypt signed by id: %s %s", signerKP.Id(), signerKP.UserName())
	f, err := os.Open(me.tempFile)
	util.Check(err)
	decrypter := SymDecrypterCreate(f, me.tempKey)
	return decrypter.Decrypt()
}

func (me *Decrypter) DecryptBytes() []byte {
	r := me.Decrypt()
	defer r.Close()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	return ret
}

func (me *Decrypter) DecryptString() string {
	return string(me.DecryptBytes())
}
