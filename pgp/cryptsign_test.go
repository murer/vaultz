package pgp

import (
	"strings"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	signed := SignString("mymsg", maria)

	decrypter := CreateDecrypter(strings.NewReader(signed)).Signers(KeyRingCreate(maria))
	defer decrypter.Close()
	assert.Equal(t, "mymsg", util.ReadAllString(decrypter.Start()))
}
