package pgp

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/murer/vaultz/util"
	"github.com/stretchr/testify/assert"
)

func TestPocSequential(t *testing.T) {
	s := SymKeyGenerate()
	kp := KeyGenerate("me", "me@sample.com")

	buf := new(bytes.Buffer)
	func() {
		encrypter := EncypterCreate(buf, kp, KeyRingCreate(kp))
		defer encrypter.Close()
		writer := encrypter.Encrypt()
		writer.Write([]byte("first"))
	}()
	buf.Write([]byte("\n"))

	func() {
		encrypter := SymEncypterCreate(buf, s)
		defer encrypter.Close()
		writer := encrypter.Encrypt()
		writer.Write([]byte("second"))
	}()

	func() {
		decrypter := DecrypterCreate(buf, KeyRingCreate(kp), KeyRingCreate())
		defer decrypter.Close()
		reader := decrypter.Decrypt()
		data, err := ioutil.ReadAll(reader)
		util.Check(err)
		assert.Equal(t, "first", string(data))
	}()

	// d, err := ioutil.ReadAll(buf)
	// util.Check(err)
	// log.Printf("rest: %s", string(d))

	func() {
		decrypter := SymDecrypterCreate(buf, s)
		defer decrypter.Close()
		reader := decrypter.Decrypt()
		data, err := ioutil.ReadAll(reader)
		util.Check(err)
		assert.Equal(t, "second", string(data))
	}()

}