package pgp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyIdCollision(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")

	assert.Equal(t, 2, KeyRingCreate(maria, bob).Size())

	assert.Panics(t, func() {
		assert.Equal(t, 2, KeyRingCreate(maria, maria).Size())
	})
}
