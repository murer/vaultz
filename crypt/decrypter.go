package crypt

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func DecrypterCreate(plain io.Reader, writers *KeyRing, readers *KeyRing) *Decrypter {
	return &Decrypter{plain: plain, writers: writers, readers: readers}
}

type Decrypter struct {
	io.Closer
	plain   io.Reader
	writers *KeyRing
	readers *KeyRing

	msg      *openpgp.MessageDetails
	tempKey  *SymKey
	tempFile string
}

func (me *Decrypter) Close() error {
	if me.tempFile != "" {
		log.Printf("Decrypter closing, delete: %s", me.tempFile)
		os.Remove(me.tempFile)
	}
	return nil
}

func (me *Decrypter) UnsafeDecrypt() io.Reader {
	ar, err := armor.Decode(me.plain)
	util.Check(err)
	keys := append(me.readers.toPgpEntityList(), me.writers.toPgpEntityList()...)
	msg, err := openpgp.ReadMessage(ar.Body, keys, nil, nil)
	util.Check(err)
	me.msg = msg
	decKP := keyFromEntity(me.msg.DecryptedWith.Entity)
	log.Printf("Decrypt with: %s %s", decKP.Id(), decKP.UserName())
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
	// TODO: must be encrypted
	unsafe := me.UnsafeDecrypt()
	f, err := ioutil.TempFile(os.TempDir(), "vaultz-decrypt-*.tmp")
	util.Check(err)
	defer f.Close()
	me.tempKey = SymKeyGenerate()
	encrypter := SymEncypterCreate(f, me.tempKey)
	defer encrypter.Close()
	total, err := io.Copy(encrypter.Encrypt(), unsafe)
	util.Check(err)
	log.Printf("Decrypt to %s, total: %d", f.Name(), total)
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
