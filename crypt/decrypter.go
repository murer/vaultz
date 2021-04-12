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

func DecrypterCreate(plain io.Reader, writers *KeyRing) *Decrypter {
	return &Decrypter{plain: plain, writers: writers}
}

type Decrypter struct {
	plain   io.Reader
	writers *KeyRing

	msg      *openpgp.MessageDetails
	tempFile string
}

func (me *Decrypter) UnsafeDecrypt() io.Reader {
	ar, err := armor.Decode(me.plain)
	util.Check(err)
	msg, err := openpgp.ReadMessage(ar.Body, me.writers.toPgpEntityList(), nil, nil)
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
	total, err := io.Copy(f, unsafe)
	util.Check(err)
	log.Printf("Decrypt to %s, total: %d", f.Name(), total)
	me.tempFile = f.Name()
}

func (me *Decrypter) Decrypt() io.ReadCloser {
	me.decryptToTemp()
	if me.msg.Signature == nil {
		msg := fmt.Sprintf("bad sign: %X", me.msg.SignedByKeyId)
		log.Panic(msg)
	}
	signerKP := keyFromEntity(me.msg.SignedBy.Entity)
	log.Printf("Decrypt signed by id: %s %s", signerKP.Id(), signerKP.UserName())
	ret, err := os.Open(me.tempFile)
	util.Check(err)
	return ret
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
