package crypt

import (
	"github.com/murer/vaultz/crypt/util"

	"golang.org/x/crypto/openpgp"
)

func KeyGen(name string, email string) *openpgp.Entity {
	var e *openpgp.Entity
	e, err := openpgp.NewEntity(name, name, email, nil)
	util.Check(err)
	return e
}
