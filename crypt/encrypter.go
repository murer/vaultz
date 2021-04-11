package crypt

import (
	"bytes"
	"io"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type encrypter struct {
	io.WriteCloser
	armor  io.WriteCloser
	writer io.WriteCloser
}

func (me *encrypter) Write(p []byte) (n int, err error) {
	return me.writer.Write(p)
}

func (me *encrypter) Close() error {
	we := me.writer.Close()
	ae := me.armor.Close()
	if we != nil {
		return we
	}
	return ae
}

func Encrypt(w io.Writer, ring *KeyRing) io.WriteCloser {
	wa, err := armor.Encode(w, "PGP MESSAGE", nil)
	util.Check(err)
	ew, err := openpgp.Encrypt(wa, ring.toPgpEntityList(), ring.first().pgpkey, nil, nil)
	util.Check(err)
	return &encrypter{armor: wa, writer: ew}
}

func EncryptBytes(plain []byte, ring *KeyRing) string {
	buf := new(bytes.Buffer)
	w := Encrypt(buf, ring)
	defer w.Close()
	w.Write(plain)
	util.Check(w.Close())
	return buf.String()
}

func EncryptString(plain string, ring *KeyRing) string {
	return EncryptBytes([]byte(plain), ring)
}
