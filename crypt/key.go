package crypt

import (
	"bytes"
	"crypto"
	"log"
	"strings"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"

	"golang.org/x/crypto/openpgp"
)

func KeyGenerate(name string, email string) *KeyPair {
	config := &packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		RSABits:                1024,
	}
	pgpkey, err := openpgp.NewEntity(name, name, email, config)
	util.Check(err)
	ret := &KeyPair{pgpkey: pgpkey}
	log.Printf("KeyGenerate: %s %s", ret.Id(), ret.UserName())
	return ret
}

func KeyImport(encodedKey string) *KeyPair {
	lst, err := openpgp.ReadArmoredKeyRing(strings.NewReader(encodedKey))
	util.Check(err)
	ret := &KeyPair{pgpkey: lst[0]}
	log.Printf("KeyImport: %s %s", ret.Id(), ret.UserName())
	return ret
}

type KeyPair struct {
	pgpkey *openpgp.Entity
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
	if me.pgpkey.PrivateKey == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	a, err := armor.Encode(buf, openpgp.PrivateKeyType, nil)
	util.Check(err)
	defer a.Close()
	me.pgpkey.SerializePrivate(a, nil)
	a.Close()
	return buf.String()
}

func (me *KeyPair) Id() string {
	return me.pgpkey.PrimaryKey.KeyIdString()
}

func (me *KeyPair) UserName() string {
	for _, k := range me.pgpkey.Identities {
		return k.UserId.Name
	}
	return ""
}

func (me *KeyPair) UserEmail() string {
	for _, k := range me.pgpkey.Identities {
		return k.UserId.Email
	}
	return ""
}

func (me *KeyPair) PubOnly() *KeyPair {
	return KeyImport(me.ExportPub())
}
