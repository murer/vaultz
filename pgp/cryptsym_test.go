package pgp

import (
	"strings"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestSymKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	k1 := SymKeyGenerate()
	k2 := SymKeyImport(k1.Export())
	assert.Equal(t, k1.key, k2.key)

	ciphered1 := SymEncryptString("mymsg", k1)
	ciphered2 := SymEncryptString("mymsg", k1)

	assert.NotEqual(t, ciphered1, ciphered2)

	// decrypter := SymDecrypterCreate(strings.NewReader(ciphered1), k1)
	// defer decrypter.Close()
	// assert.Equal(t, "mymsg", decrypter.DecryptString())
	// decrypter = SymDecrypterCreate(strings.NewReader(ciphered2), k1)
	// defer decrypter.Close()
	// assert.Equal(t, "mymsg", decrypter.DecryptString())

	decrypter := CreateDecrypter(strings.NewReader(ciphered1))
	defer decrypter.Close()
	assert.Equal(t, "mymsg", util.ReadAllString(decrypter.Symmetric(k1).Open()))
	decrypter = CreateDecrypter(strings.NewReader(ciphered2))
	defer decrypter.Close()
	assert.Equal(t, "mymsg", util.ReadAllString(decrypter.Symmetric(k1).Open()))

}
