package crypt

import (
	"io"
	"io/ioutil"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func DecrypterCreate(plain io.Reader, ring *KeyRing) *Decrypter {
	return &Decrypter{plain: plain, ring: ring}
}

type Decrypter struct {
	// io.ReadCloser
	// reader io.Reader

	plain io.Reader
	ring  *KeyRing

	msg *openpgp.MessageDetails
}

func (me *Decrypter) UnsafeDecrypt() io.Reader {
	ar, err := armor.Decode(me.plain)
	util.Check(err)
	msg, err := openpgp.ReadMessage(ar.Body, me.ring.toPgpEntityList(), nil, nil)
	util.Check(err)
	me.msg = msg
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

func (me *Decrypter) Decrypt() io.Reader {
	ar, err := armor.Decode(me.plain)
	util.Check(err)
	msg, err := openpgp.ReadMessage(ar.Body, me.ring.toPgpEntityList(), nil, nil)
	util.Check(err)
	me.msg = msg
	return me.msg.UnverifiedBody
}

func (me *Decrypter) DecryptBytes() []byte {
	r := me.Decrypt()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	return ret
}

func (me *Decrypter) DecryptString() string {
	return string(me.DecryptBytes())
}
