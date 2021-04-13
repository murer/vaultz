package pgp

import (
	"bytes"
	"io"
	"strings"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp/armor"
)

func ArmorEncode(writer io.Writer, blockType string) io.WriteCloser {
	ar, err := armor.Encode(writer, blockType, nil)
	util.Check(err)
	return ar
}

func ArmorDecode(reader io.Reader) io.Reader {
	ar, err := armor.Decode(reader)
	util.Check(err)
	return ar.Body
}

func ArmorEncodeBytes(data []byte, blockType string) string {
	buf := new(bytes.Buffer)
	func() {
		writer := ArmorEncode(buf, blockType)
		defer writer.Close()
		writer.Write(data)
	}()
	return buf.String()
}

func ArmorEncodeString(data string, blockType string) string {
	return ArmorEncodeBytes([]byte(data), blockType)
}

func ArmorDecodeBytes(data string, blockType string) []byte {
	reader := ArmorDecode(strings.NewReader(data))
	return util.ReadAll(reader)
}

func ArmorDecodeString(data string, blockType string) string {
	return string(ArmorDecodeBytes(data, blockType))
}
