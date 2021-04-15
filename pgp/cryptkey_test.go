package pgp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 16, len(kp.Id()))

	fmt.Printf("id: %s\n", kp.Id())

	pubkp := KeyImport(kp.ExportPubArmored())
	assert.Equal(t, "test", pubkp.UserName())
	assert.Equal(t, "test@sample.com", pubkp.UserEmail())
	assert.Equal(t, kp.ExportPubArmored(), pubkp.ExportPubArmored())
	assert.Empty(t, pubkp.ExportPrivArmored())
	assert.Equal(t, kp.Id(), pubkp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	privkp := KeyImport(kp.ExportPrivArmored())
	assert.Equal(t, "test", privkp.UserName())
	assert.Equal(t, "test@sample.com", privkp.UserEmail())
	assert.Equal(t, kp.ExportPubArmored(), privkp.ExportPubArmored())
	assert.NotEmpty(t, privkp.ExportPrivArmored())
	assert.Equal(t, kp.Id(), privkp.Id())
	fmt.Printf("id: %s\n", privkp.Id())

}
