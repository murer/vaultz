package pgp

import (
	"bytes"
	"strings"

	"github.com/murer/vaultz/util"
)

func ArmorEncodeBytes(data []byte, blockType string) string {
	buf := new(bytes.Buffer)
	func() {
		enc := CreateEncrypter(buf).Armored(blockType)
		defer enc.Close()
		enc.Start().Write(data)
	}()
	return buf.String()
}

func ArmorEncodeString(data string, blockType string) string {
	return ArmorEncodeBytes([]byte(data), blockType)
}

func ArmorDecodeBytes(data string, blockType string) []byte {
	dec := CreateDecrypter(strings.NewReader(data)).Armored(true)
	return util.ReadAll(dec.Start())
}

func ArmorDecodeString(data string, blockType string) string {
	return string(ArmorDecodeBytes(data, blockType))
}
