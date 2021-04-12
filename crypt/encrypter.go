package crypt

import (
	"bytes"
	"io"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type Encrypter struct {
	io.WriteCloser

	ciphered   io.Writer
	signer     *KeyPair
	recipients *KeyRing

	armor  io.WriteCloser
	writer io.WriteCloser
}

func EncypterCreate(ciphered io.Writer, signer *KeyPair, recipients *KeyRing) *Encrypter {
	return &Encrypter{
		ciphered:   ciphered,
		signer:     signer,
		recipients: recipients,
	}
}

func (me *Encrypter) Encrypt() io.WriteCloser {
	wa, err := armor.Encode(me.ciphered, "PGP MESSAGE", nil)
	util.Check(err)
	ew, err := openpgp.Encrypt(wa, me.recipients.toPgpEntityList(), me.signer.pgpkey, nil, nil)
	util.Check(err)
	me.armor = wa
	me.writer = ew
	return me
}

func (me *Encrypter) Write(p []byte) (n int, err error) {
	return me.writer.Write(p)
}

func (me *Encrypter) Close() error {
	we := me.writer.Close()
	ae := me.armor.Close()
	if we != nil {
		return we
	}
	return ae
}

func EncryptBytes(plain []byte, signer *KeyPair, recipients *KeyRing) string {
	buf := new(bytes.Buffer)
	encrypter := EncypterCreate(buf, signer, recipients)
	w := encrypter.Encrypt()
	defer w.Close()
	w.Write(plain)
	util.Check(w.Close())
	return buf.String()
}

func EncryptString(plain string, signer *KeyPair, ring *KeyRing) string {
	return EncryptBytes([]byte(plain), signer, ring)
}
