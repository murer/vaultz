package pgp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptWrongDecrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	recipients := KeyRingCreate(bob)
	ciphered := EncryptString("mymsg", maria, recipients)
	writers := KeyRingCreate(maria.PubOnly())
	decrypter := DecrypterCreate(bytes.NewReader(ciphered), writers, KeyRingCreate())
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
	decrypter := DecrypterCreate(bytes.NewReader(ciphered), KeyRingCreate(bob), KeyRingCreate(bob))
	defer decrypter.Close()
	assert.Panics(t, func() {
		decrypter.DecryptString()
	})
}

func TestCryptWrongWriter(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", john, recipients)
	readers := KeyRingCreate(john)
	decrypter := DecrypterCreate(bytes.NewReader(ciphered), KeyRingCreate(), readers)
	defer decrypter.Close()
	writers := KeyRingCreate(maria)
	decrypter = DecrypterCreate(bytes.NewReader(ciphered), writers, readers)
	defer decrypter.Close()
	assert.Panics(t, func() {
		decrypter.DecryptString()
	})
}

func TestCryptUnsigned(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(maria, bob, john)
	ciphered := EncryptString("mymsg", nil, recipients)
	readers := KeyRingCreate(john)
	decrypter := DecrypterCreate(bytes.NewReader(ciphered), nil, readers)
	defer decrypter.Close()
	plain := decrypter.DecryptString()
	assert.Equal(t, "mymsg", plain)
}

func TestCryptSignUncheck(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(maria, bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	readers := KeyRingCreate(john)
	decrypter := DecrypterCreate(bytes.NewReader(ciphered), nil, readers)
	defer decrypter.Close()
	plain := decrypter.DecryptString()
	assert.Equal(t, "mymsg", plain)
}

func TestCrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(maria, bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	// fmt.Println(maria.ExportPriv())
	// fmt.Println(ArmorEncodeBytes(ciphered, "PGP MESSAGE"))
	readers := KeyRingCreate(john)
	decrypter := DecrypterCreate(bytes.NewReader(ciphered), KeyRingCreate(), readers)
	defer decrypter.Close()
	unsafePlain := decrypter.UnsafeDecryptString()
	assert.Equal(t, "mymsg", unsafePlain)
	writers := KeyRingCreate(maria)
	decrypter = DecrypterCreate(bytes.NewReader(ciphered), writers, readers)
	defer decrypter.Close()
	plain := decrypter.DecryptString()
	assert.Equal(t, "mymsg", plain)
}
