package pgp

import (
	"bytes"
	"strings"

	"github.com/murer/vaultz/util"
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

func DecryptBytes(ciphered []byte, signers *KeyRing, recipients *KeyRing) []byte {
	dec := CreateDecrypter(bytes.NewBuffer(ciphered)).Signers(signers).Decrypt(recipients)
	defer dec.Close()
	return util.ReadAll(dec.Start())
}

func DecryptString(ciphered []byte, signers *KeyRing, recipients *KeyRing) string {
	return string(DecryptBytes(ciphered, signers, recipients))
}

func SymEncryptBytes(plain []byte, key *SymKey) []byte {
	buf := new(bytes.Buffer)
	func() {
		encrypter := CreateEncrypter(buf).Symmetric(key)
		defer encrypter.Close()
		encrypter.Start().Write(plain)
	}()
	return buf.Bytes()
}

func SymEncryptString(plain string, key *SymKey) []byte {
	return SymEncryptBytes([]byte(plain), key)
}

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

func SignBytes(plain []byte, signer *KeyPair) string {
	buf := new(bytes.Buffer)
	func() {
		enc := CreateEncrypter(buf).Sign(signer)
		defer enc.Close()
		enc.Start().Write(plain)
	}()
	return buf.String()
}

func SignString(plain string, signer *KeyPair) string {
	return SignBytes([]byte(plain), signer)
}
