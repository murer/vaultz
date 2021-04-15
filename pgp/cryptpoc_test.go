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
		encrypter := CreateEncrypter(buf).Sign(kp).Encrypt(KeyRingCreate(kp))
		defer encrypter.Close()
		writer := encrypter.Start()
		writer.Write([]byte("first"))
	}()

	func() {
		encrypter := SymEncypterCreate(buf, s)
		defer encrypter.Close()
		writer := encrypter.Encrypt()
		writer.Write([]byte("second"))
	}()

	func() {
		decrypter := CreateDecrypter(buf).Signers(KeyRingCreate(kp)).Decrypt(KeyRingCreate())
		defer decrypter.Close()
		reader := decrypter.Start()
		data, err := ioutil.ReadAll(reader)
		util.Check(err)
		assert.Equal(t, "first", string(data))
	}()

	// d, err := ioutil.ReadAll(buf)
	// util.Check(err)
	// log.Printf("rest: %s", string(d))

	func() {
		decrypter := CreateDecrypter(buf).Symmetric(s)
		defer decrypter.Close()
		reader := decrypter.Start()
		data, err := ioutil.ReadAll(reader)
		util.Check(err)
		assert.Equal(t, "second", string(data))
	}()

}
