package pgp

import (
	"strings"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestArmor(t *testing.T) {

	armored := ArmorEncodeString("mymsg", "TEST ONLY")

	reader := CreateDecrypter(strings.NewReader(armored)).Armor(true).Open()
	defer reader.Close()
	assert.Equal(t, "mymsg", string(util.ReadAll(reader)))
}
