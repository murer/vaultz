package crypt

import (
	"github.com/murer/vaultz/crypt/util"

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
