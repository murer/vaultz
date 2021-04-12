package pgp

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func SymDecrypterCreate(plain io.Reader, key *SymKey) *SymDecrypter {
	return &SymDecrypter{plain: plain, key: key}
}

type SymDecrypter struct {
	io.ReadCloser
	plain io.Reader
	key   *SymKey

	msg *openpgp.MessageDetails
}

func (me *SymDecrypter) Read(p []byte) (n int, err error) {
	return me.msg.UnverifiedBody.Read(p)
}

func (me *SymDecrypter) Close() error {
	return nil
}

func (me *SymDecrypter) Decrypt() io.ReadCloser {
	ar, err := armor.Decode(me.plain)
	util.Check(err)
	msg, err := openpgp.ReadMessage(ar.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return me.key.key, nil
	}, nil)
	util.Check(err)
	me.msg = msg
	log.Printf("SymDecrypt with key size: %d", me.key.Size())
	return me
}

func (me *SymDecrypter) DecryptBytes() []byte {
	r := me.Decrypt()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	return ret
}

func (me *SymDecrypter) DecryptString() string {
	return string(me.DecryptBytes())
}
