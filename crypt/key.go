package crypt

import (
	"bytes"
	"strings"

	"golang.org/x/crypto/openpgp/packet"

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

func (me *KeyPair) Import(encodedKey string) *KeyPair {
	blk, err := armor.Decode(strings.NewReader(encodedKey))
	util.Check(err)
	reader := packet.NewReader(blk.Body)
	pgpkey, err := openpgp.ReadEntity(reader)
	util.Check(err)
	me.pgpkey = pgpkey
	// for {
	// 	fmt.Println("Reading")
	// 	pkt, err := reader.Next()
	// 	util.Check(err)
	// 	fmt.Printf("AA %#v \n", pkt)
	// }

	// key, ok := pkt.(*packet.PrivateKey)
	// if !ok {
	// 	util.Fatal("Not private key")
	// }
	return me
}
