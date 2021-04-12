package crypt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := KeyGenerate("test", "test@sample.com")
	assert.Equal(t, "test", kp.UserName())
	assert.Equal(t, "test@sample.com", kp.UserEmail())
	assert.NotEmpty(t, kp.ExportPub())
	assert.NotEmpty(t, kp.ExportPriv())
	assert.NotEmpty(t, kp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	pubkp := KeyImport(kp.ExportPub())
	assert.Equal(t, "test", pubkp.UserName())
	assert.Equal(t, "test@sample.com", pubkp.UserEmail())
	assert.Equal(t, kp.ExportPub(), pubkp.ExportPub())
	assert.Empty(t, pubkp.ExportPriv())
	assert.Equal(t, kp.Id(), pubkp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	privkp := KeyImport(kp.ExportPriv())
	assert.Equal(t, "test", privkp.UserName())
	assert.Equal(t, "test@sample.com", privkp.UserEmail())
	assert.Equal(t, kp.ExportPub(), privkp.ExportPub())
	assert.NotEmpty(t, privkp.ExportPriv())
	assert.Equal(t, kp.Id(), privkp.Id())
	fmt.Printf("id: %s\n", privkp.Id())
}

func TestCrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	// fmt.Println(maria.ExportPub())
	// fmt.Println(bob.ExportPriv())
	ring := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", maria, ring)
	// fmt.Println(ciphered)
	ring = KeyRingCreate(maria.PubOnly(), john)
	decrypter := DecrypterCreate(strings.NewReader(ciphered), ring)
	unsafePlain := decrypter.UnsafeDecryptString()
	assert.Equal(t, "mymsg", unsafePlain)
	decrypter = DecrypterCreate(strings.NewReader(ciphered), ring)
	plain := decrypter.DecryptString()
	assert.Equal(t, "mymsg", plain)
}
