package pgp

import (
	"strings"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestArmor(t *testing.T) {

	armored := ArmorEncodeString("mymsg", "TEST ONLY")

	decrypter := CreateDecrypter(strings.NewReader(armored))
	defer decrypter.Close()
	reader := decrypter.Armor(true).Open()
	assert.Equal(t, "mymsg", string(util.ReadAll(reader)))
}
