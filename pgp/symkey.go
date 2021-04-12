package pgp

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"log"
	"strings"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp/armor"
)

func SymKeyGenerate() *SymKey {
	size := 32
	b := make([]byte, size)
	s, err := rand.Read(b)
	util.Check(err)
	if s != size {
		log.Panicf("it must be %d, but was %d", size, s)
	}
	ret := &SymKey{key: b}
	log.Printf("SymKeyGenerate, size: %d", ret.Size())
	return ret
}

func SymKeyImport(encodedKey string) *SymKey {
	ar, err := armor.Decode(strings.NewReader(encodedKey))
	util.Check(err)
	key, err := ioutil.ReadAll(ar.Body)
	util.Check(err)
	ret := &SymKey{key: key}
	log.Printf("KeyImport, size: %d", ret.Size())
	return ret
}

type SymKey struct {
	key []byte
}

func (me *SymKey) Size() int {
	return len(me.key)
}

func (me *SymKey) Export() string {
	buf := new(bytes.Buffer)
	a, err := armor.Encode(buf, "VAULTZ SYMMETRIC KEY", nil)
	util.Check(err)
	defer a.Close()
	a.Write(me.key)
	a.Close()
	return buf.String()
}
