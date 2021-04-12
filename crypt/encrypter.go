package crypt

import (
	"bytes"
	"io"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type encrypterWriter struct {
	io.WriteCloser
	armor  io.WriteCloser
	writer io.WriteCloser
}

func (me *encrypterWriter) Write(p []byte) (n int, err error) {
	return me.writer.Write(p)
}

func (me *encrypterWriter) Close() error {
	we := me.writer.Close()
	ae := me.armor.Close()
	if we != nil {
		return we
	}
	return ae
}

func Encrypt(w io.Writer, signer *KeyPair, ring *KeyRing) io.WriteCloser {
	wa, err := armor.Encode(w, "PGP MESSAGE", nil)
	util.Check(err)
	ew, err := openpgp.Encrypt(wa, ring.toPgpEntityList(), signer.pgpkey, nil, nil)
	util.Check(err)
	return &encrypterWriter{armor: wa, writer: ew}
}

func EncryptBytes(plain []byte, signer *KeyPair, ring *KeyRing) string {
	buf := new(bytes.Buffer)
	w := Encrypt(buf, signer, ring)
	defer w.Close()
	w.Write(plain)
	util.Check(w.Close())
	return buf.String()
}

func EncryptString(plain string, signer *KeyPair, ring *KeyRing) string {
	return EncryptBytes([]byte(plain), signer, ring)
}
