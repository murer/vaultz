package pgp

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/openpgp"
)

func TestKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := KeyGenerate("test", "test@sample.com")
	assert.Equal(t, "test", kp.UserName())
	assert.Equal(t, "test@sample.com", kp.UserEmail())
	assert.NotEmpty(t, kp.ExportPubArmored())
	assert.Regexp(t, "-----END PGP PUBLIC KEY BLOCK-----$", kp.ExportPubArmored())
	assert.NotEmpty(t, kp.ExportPrivArmored())
	assert.Regexp(t, "-----END PGP PRIVATE KEY BLOCK-----$", kp.ExportPrivArmored())
	assert.NotEmpty(t, kp.Id())
	assert.Equal(t, 16, len(kp.IdString()))

	assert.Equal(t, fmt.Sprintf("%016X", kp.Id()), kp.IdString())

	pubkp := KeyImport(kp.ExportPubArmored())
	assert.Equal(t, "test", pubkp.UserName())
	assert.Equal(t, "test@sample.com", pubkp.UserEmail())
	assert.Equal(t, kp.ExportPubArmored(), pubkp.ExportPubArmored())
	assert.Empty(t, pubkp.ExportPrivArmored())
	assert.Equal(t, kp.Id(), pubkp.Id())
	fmt.Printf("id: %X\n", kp.Id())

	privkp := KeyImport(kp.ExportPrivArmored())
	assert.Equal(t, "test", privkp.UserName())
	assert.Equal(t, "test@sample.com", privkp.UserEmail())
	assert.Equal(t, kp.ExportPubArmored(), privkp.ExportPubArmored())
	assert.NotEmpty(t, privkp.ExportPrivArmored())
	assert.Equal(t, kp.Id(), privkp.Id())
	fmt.Printf("id: %X\n", privkp.Id())

}

func TestKeyRingSerialization(t *testing.T) {
	a := KeyGenerate("a", "a@sample.com")
	b := KeyGenerate("b", "b@sample.com")
	c := KeyGenerate("c", "c@sample.com")
	ring := KeyRingCreate(a, b, c)

	buf := new(bytes.Buffer)
	for _, key := range ring.kps {
		data := key.ExportPubBinary()
		buf.Write(data)
	}

	lst, err := openpgp.ReadKeyRing(buf)
	util.Check(err)
	log.Printf("ents %#v", lst)
	assert.Equal(t, 3, len(lst))

	nring := KeyRingCreate().fromPgpEntityList(lst...)
	log.Printf("ring %#v", nring)
	assert.Equal(t, a.ExportPubArmored(), nring.Get(a.Id()).ExportPubArmored())
	assert.Equal(t, b.ExportPubArmored(), nring.Get(b.Id()).ExportPubArmored())
	assert.Equal(t, c.ExportPubArmored(), nring.Get(c.Id()).ExportPubArmored())
	assert.Equal(t, 3, nring.Size())
}
