package pgp

import (
	"bytes"
)

func EncryptBytes(plain []byte, signer *KeyPair, recipients *KeyRing) []byte {
	buf := new(bytes.Buffer)
	func() {
		encrypter := CreateEncrypter(buf).Sign(signer).Encrypt(recipients)
		defer encrypter.Close()
		encrypter.Start().Write(plain)
	}()
	return buf.Bytes()
}

func EncryptString(plain string, signer *KeyPair, ring *KeyRing) []byte {
	return EncryptBytes([]byte(plain), signer, ring)
}
