package crypt

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	k1 := SymKeyGenerate()
	k2 := SymKeyImport(k1.Export())
	assert.Equal(t, k1.key, k2.key)
	log.Printf("XXX %X", k1.key)

	ciphered1 := SymEncryptString("mymsg", k1)
	fmt.Println(ciphered1)
	ciphered2 := SymEncryptString("mymsg", k1)
	fmt.Println(ciphered2)

	assert.NotEqual(t, ciphered1, ciphered2)

	decrypter := SymDecrypterCreate(strings.NewReader(ciphered1), k1)
	defer decrypter.Close()
	assert.Equal(t, "mymsg", decrypter.DecryptString())
	decrypter = SymDecrypterCreate(strings.NewReader(ciphered2), k1)
	defer decrypter.Close()
	assert.Equal(t, "mymsg", decrypter.DecryptString())
}
