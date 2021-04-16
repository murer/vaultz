package pgp

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Printf("RSABits is set to 1024 in test")
	Config.RSABits = 1024
}

func TestPacketDesc(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")

	log.Printf("Key bob: %X", bob.Id())
	log.Printf("Key john: %X", john.Id())

	recipients := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	assert.Equal(t, "mymsg", DecryptString(ciphered, KeyRingCreate(maria), recipients))

	assert.Panics(t, func() {
		assert.Equal(t, "mymsg", DecryptString(ciphered, nil, KeyRingCreate()))
	})

	_, err := TryToDecryptString(ciphered, nil, KeyRingCreate())
	assert.Equal(t, ErrKeyIncorrect, err)

	_, err = TryToDecryptString(ciphered, KeyRingCreate(maria), KeyRingCreate(maria))
	assert.Equal(t, ErrKeyIncorrect, err)

	_, err = TryToDecryptString(ciphered, KeyRingCreate(maria), KeyRingCreate())
	assert.Equal(t, ErrKeyIncorrect, err)
}
