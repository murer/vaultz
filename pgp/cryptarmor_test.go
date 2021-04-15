package pgp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestArmor(t *testing.T) {

	buf := new(bytes.Buffer)
	func() {
		encrypter := CreateEncrypter(buf).Armored("TEST ONLY")
		defer encrypter.Close()
		encrypter.Start().Write([]byte("mymsg"))
	}()
	decrypter := CreateDecrypter(strings.NewReader(buf.String())).Armor(true)
	defer decrypter.Close()
	assert.Equal(t, "mymsg", string(util.ReadAll(decrypter.Start())))
}
