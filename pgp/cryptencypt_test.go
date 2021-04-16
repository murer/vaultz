package pgp

import (
	"bytes"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestCryptWrongDecrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("jphn", "john@sample.com")
	recipients := KeyRingCreate(bob)
	ciphered := EncryptString("mymsg", maria, recipients)
	writers := KeyRingCreate(maria.PubOnly())
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(writers).Recipients(KeyRingCreate(john))
	defer decrypter.Close()
	assert.Panics(t, func() {
		util.ReadAll(decrypter.Start())
	})
}

func TestCryptWrongSign(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	recipients := KeyRingCreate(bob)
	ciphered := EncryptString("mymsg", maria, recipients)
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(KeyRingCreate(bob)).Recipients(KeyRingCreate(bob))
	defer decrypter.Close()
	assert.Panics(t, func() {
		util.ReadAll(decrypter.Start())
	})
}

func TestCryptWrongWriter(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", john, recipients)
	readers := KeyRingCreate(john)
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(KeyRingCreate()).Recipients(readers)
	defer decrypter.Close()
	assert.Panics(t, func() {
		util.ReadAll(decrypter.Start())
	})
	writers := KeyRingCreate(maria)
	decrypter = CreateDecrypter(bytes.NewReader(ciphered)).Signers(writers).Recipients(readers)
	defer decrypter.Close()
	assert.Panics(t, func() {
		util.ReadAll(decrypter.Start())
	})
}

func TestCryptUnsigned(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(maria, bob, john)
	ciphered := EncryptString("mymsg", nil, recipients)
	readers := KeyRingCreate(john)
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(nil).Recipients(readers)
	defer decrypter.Close()
	plain := util.ReadAllString(decrypter.Start())
	assert.Equal(t, "mymsg", plain)
}

func TestCryptSignUncheck(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(maria, bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	readers := KeyRingCreate(john)
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(nil).Recipients(readers)
	defer decrypter.Close()
	plain := util.ReadAllString(decrypter.Start())
	assert.Equal(t, "mymsg", plain)
}

func TestCryptLarge(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	ring := KeyRingCreate(maria)
	ciphered := EncryptString("01234567890123456789", maria, ring)
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(nil).Recipients(ring).MaxTempMemory(10)
	defer decrypter.Close()
	unsafePlain := util.ReadAllString(decrypter.Start())
	assert.Equal(t, "01234567890123456789", unsafePlain)
	decrypter = CreateDecrypter(bytes.NewReader(ciphered)).Signers(ring).Recipients(ring).MaxTempMemory(10)
	defer decrypter.Close()
	plain := util.ReadAllString(decrypter.Start())
	assert.Equal(t, "01234567890123456789", plain)
}

func TestCrypt(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")
	recipients := KeyRingCreate(maria, bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	// fmt.Println(maria.ExportPrivArmored())
	// fmt.Println(ArmorEncodeBytes(ciphered, "PGP MESSAGE"))
	readers := KeyRingCreate(john)
	decrypter := CreateDecrypter(bytes.NewReader(ciphered)).Signers(nil).Recipients(readers)
	defer decrypter.Close()
	unsafePlain := util.ReadAllString(decrypter.Start())
	assert.Equal(t, "mymsg", unsafePlain)
	writers := KeyRingCreate(maria)
	decrypter = CreateDecrypter(bytes.NewReader(ciphered)).Signers(writers).Recipients(readers)
	defer decrypter.Close()
	plain := util.ReadAllString(decrypter.Start())
	assert.Equal(t, "mymsg", plain)
}
