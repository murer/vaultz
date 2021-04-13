package pgp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	signed := SignString("mymsg", maria)
	fmt.Println(signed)

	decrypter := VerifierCreate(strings.NewReader(signed), KeyRingCreate(maria))
	defer decrypter.Close()
	assert.Equal(t, "mymsg", decrypter.DecryptString())
}
