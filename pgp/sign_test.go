package pgp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	signed := SignString("mymsg", maria)

	decrypter := VerifierCreate(strings.NewReader(signed), KeyRingCreate(maria))
	defer decrypter.Close()
	assert.Equal(t, "mymsg", decrypter.DecryptString())
}
