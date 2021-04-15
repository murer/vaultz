package pgp

import (
	"bytes"
	"crypto"
	"log"
	"strings"

	"github.com/murer/vaultz/util"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
)

func KeyGenerate(name string, email string) *KeyPair {
	pgpkey, err := openpgp.NewEntity(name, name, email, Config)
	util.Check(err)
	ret := &KeyPair{pgpkey: pgpkey}
	log.Printf("KeyGenerate: %s %s", ret.Id(), ret.UserName())
	return ret
}

func KeyImport(encodedKey string) *KeyPair {
	lst, err := openpgp.ReadArmoredKeyRing(strings.NewReader(encodedKey))
	util.Check(err)
	ret := &KeyPair{pgpkey: lst[0]}
	for k, v := range ret.pgpkey.Identities {
		log.Printf("X: %s = '%#v'\n", k, v)
		v.SelfSignature.PreferredSymmetric = []uint8{uint8(packet.CipherAES256)}
		id, ok := s2k.HashToHashId(crypto.SHA256)
		util.Assert(!ok, "hash not found")
		v.SelfSignature.PreferredHash = []uint8{id}
	}
	return ret
}

func keyFromEntity(entity *openpgp.Entity) *KeyPair {
	return &KeyPair{pgpkey: entity}
}

type KeyPair struct {
	pgpkey *openpgp.Entity
}

func (me *KeyPair) ExportPubBinary() []byte {
	buf := new(bytes.Buffer)
	util.Check(me.pgpkey.Serialize(buf))
	return buf.Bytes()
}

func (me *KeyPair) ExportPrivBinary() []byte {
	buf := new(bytes.Buffer)
	util.Check(me.pgpkey.SerializePrivate(buf, Config))
	return buf.Bytes()
}

func (me *KeyPair) ExportPubArmored() string {
	return ArmorEncodeBytes([]byte(me.ExportPubBinary()), openpgp.PublicKeyType)
}

func (me *KeyPair) ExportPrivArmored() string {
	if me.pgpkey.PrivateKey == nil {
		return ""
	}
	return ArmorEncodeBytes([]byte(me.ExportPrivBinary()), openpgp.PrivateKeyType)
}

func (me *KeyPair) Id() string {
	return me.pgpkey.PrimaryKey.KeyIdString()
}

func (me *KeyPair) IdBinary() uint64 {
	return me.pgpkey.PrimaryKey.KeyId
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
	return KeyImport(me.ExportPubArmored())
}
