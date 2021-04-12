package crypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	k1 := SymKeyGenerate()
	k2 := SymKeyImport(k1.Export())
	assert.Equal(t, k1.key, k2.key)

	cyphered := SymEncryptString("mymsg", k1)
	fmt.Println(cyphered)

}
