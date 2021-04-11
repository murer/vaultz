package crypt

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type decrypter struct {
	io.ReadCloser
	reader io.Reader
}

func (me *decrypter) Read(p []byte) (n int, err error) {
	return me.reader.Read(p)
}

func (me *decrypter) Close() error {
	return nil
}

func Decrypt(r io.Reader, ring *KeyRing) io.ReadCloser {
	ar, err := armor.Decode(r)
	util.Check(err)
	in, err := openpgp.ReadMessage(ar.Body, ring.toPgpEntityList(), nil, nil)
	util.Check(err)
	return &decrypter{reader: in.UnverifiedBody}
}

func DecryptBytes(cipher string, ring *KeyRing) []byte {
	r := Decrypt(strings.NewReader(cipher), ring)
	defer r.Close()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	util.Check(r.Close())
	return ret
}

func DecryptString(cipher string, ring *KeyRing) string {
	return string(EncryptBytes([]byte(cipher), ring))
}
