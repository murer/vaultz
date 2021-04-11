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
	armor  io.WriteCloser
	writer io.WriteCloser
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

func Encrypt(w io.Writer, ring *KeyRing) io.WriteCloser {
	wa, err := armor.Encode(w, "PGP MESSAGE", nil)
	util.Check(err)
	ew, err := openpgp.Encrypt(wa, ring.toPgpEntityList(), ring.first().pgpkey, nil, nil)
	util.Check(err)
	return &Encrypter{armor: wa, writer: ew}
}

func EncryptString(plain string, ring *KeyRing) string {
	buf := new(bytes.Buffer)
	w := Encrypt(buf, ring)
	defer w.Close()
	w.Write([]byte(plain))
	util.Check(w.Close())
	return buf.String()
}
