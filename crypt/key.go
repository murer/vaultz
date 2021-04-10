package crypt

import (
	"bytes"

	"github.com/murer/vaultz/crypt/util"
	"golang.org/x/crypto/openpgp/armor"

	"golang.org/x/crypto/openpgp"
)

type KeyPair struct {
	pgpkey *openpgp.Entity
}

func (me *KeyPair) Generate(name string, email string) *KeyPair {
	var pgpkey *openpgp.Entity
	pgpkey, err := openpgp.NewEntity(name, name, email, nil)
	util.Check(err)
	me.pgpkey = pgpkey
	return me
}

func (me *KeyPair) ExportPub() string {
	buf := new(bytes.Buffer)
	a, err := armor.Encode(buf, openpgp.PublicKeyType, nil)
	util.Check(err)
	defer a.Close()
	util.Check(me.pgpkey.Serialize(a))
	a.Close()
	return buf.String()
}

func (me *KeyPair) ExportPriv() string {
	buf := new(bytes.Buffer)
	a, err := armor.Encode(buf, openpgp.PrivateKeyType, nil)
	util.Check(err)
	defer a.Close()
	me.pgpkey.SerializePrivate(a, nil)
	a.Close()
	return buf.String()
}
