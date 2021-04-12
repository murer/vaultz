package pgp

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

func TestCryptWrongDecrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	recipients := KeyRingCreate(bob)
	ciphered := EncryptString("mymsg", maria, recipients)
	writers := KeyRingCreate(maria.PubOnly())
	decrypter := DecrypterCreate(strings.NewReader(ciphered), writers, KeyRingCreate())
	defer decrypter.Close()
	assert.Panics(t, func() {
		decrypter.DecryptString()
	})
}

func TestCryptWrongSign(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	recipients := KeyRingCreate(bob)
	ciphered := EncryptString("mymsg", maria, recipients)
	decrypter := DecrypterCreate(strings.NewReader(ciphered), KeyRingCreate(bob), KeyRingCreate(bob))
	defer decrypter.Close()
	assert.Panics(t, func() {
		decrypter.DecryptString()
	})
}

func TestCrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	fmt.Println(maria.ExportPub())
	fmt.Println(bob.ExportPriv())
	recipients := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	fmt.Println(ciphered)
	readers := KeyRingCreate(john)
	decrypter := DecrypterCreate(strings.NewReader(ciphered), KeyRingCreate(), readers)
	defer decrypter.Close()
	unsafePlain := decrypter.UnsafeDecryptString()
	assert.Equal(t, "mymsg", unsafePlain)
	writers := KeyRingCreate(maria)
	decrypter = DecrypterCreate(strings.NewReader(ciphered), writers, readers)
	defer decrypter.Close()
	plain := decrypter.DecryptString()
	assert.Equal(t, "mymsg", plain)
}

func TestCryptWrongWriter(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", john, recipients)
	readers := KeyRingCreate(john)
	decrypter := DecrypterCreate(strings.NewReader(ciphered), KeyRingCreate(), readers)
	defer decrypter.Close()
	writers := KeyRingCreate(maria)
	decrypter = DecrypterCreate(strings.NewReader(ciphered), writers, readers)
	defer decrypter.Close()
	assert.Panics(t, func() {
		decrypter.DecryptString()
	})
}
