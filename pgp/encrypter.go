package pgp

import (
	"bytes"
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
)

type Encrypter struct {
	io.WriteCloser

	ciphered   io.Writer
	signer     *KeyPair
	recipients *KeyRing

	// armor  io.WriteCloser
	writer io.WriteCloser

	byteCount uint64
}

func EncypterCreate(ciphered io.Writer, signer *KeyPair, recipients *KeyRing) *Encrypter {
	return &Encrypter{
		ciphered:   ciphered,
		signer:     signer,
		recipients: recipients,
	}
}

func (me *Encrypter) Encrypt() io.WriteCloser {
	// wa, err := armor.Encode(me.ciphered, "PGP MESSAGE", nil)
	// util.Check(err)
	ew, err := openpgp.Encrypt(me.ciphered, me.recipients.toPgpEntityList(), me.signer.pgpkey, nil, nil)
	util.Check(err)
	// me.armor = wa
	me.writer = ew
	log.Printf("Encrypt start, signer: %s %s, total recipients: %d", me.signer.Id(), me.signer.UserName(), len(me.recipients.kps))
	for _, v := range me.recipients.kps {
		log.Printf("Encrypt start, recipients: %s %s", v.Id(), v.UserName())
	}
	return me
}

func (me *Encrypter) Write(p []byte) (n int, err error) {
	me.byteCount = me.byteCount + uint64(len(p))
	return me.writer.Write(p)
}

func (me *Encrypter) Close() error {
	log.Printf("Encrypt done, size: %d", me.byteCount)
	we := me.writer.Close()
	// ae := me.armor.Close()
	// if we != nil {
	return we
	// }
	// return ae
}

func _encryptBytes(plain []byte, signer *KeyPair, recipients *KeyRing) *bytes.Buffer {
	buf := new(bytes.Buffer)
	encrypter := EncypterCreate(buf, signer, recipients)
	w := encrypter.Encrypt()
	defer w.Close()
	w.Write(plain)
	return buf
}

func EncryptBytes(plain []byte, signer *KeyPair, recipients *KeyRing) []byte {
	return _encryptBytes(plain, signer, recipients).Bytes()
}

func EncryptString(plain string, signer *KeyPair, ring *KeyRing) []byte {
	return EncryptBytes([]byte(plain), signer, ring)
}
