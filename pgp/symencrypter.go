package pgp

import (
	"bytes"
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type SymEncrypter struct {
	io.WriteCloser

	ciphered io.Writer
	key      *SymKey

	writer io.WriteCloser

	byteCount uint64
}

func SymEncypterCreate(ciphered io.Writer, key *SymKey) *SymEncrypter {
	return &SymEncrypter{
		ciphered: ciphered,
		key:      key,
	}
}

func (me *SymEncrypter) Encrypt() io.WriteCloser {
	packetConfig := &packet.Config{
		DefaultCipher: packet.CipherAES256,
	}
	wa, err := armor.Encode(me.ciphered, "PGP MESSAGE", nil)
	util.Check(err)
	ew, err := openpgp.SymmetricallyEncrypt(wa, me.key.key, nil, packetConfig)
	util.Check(err)
	me.writer = ew
	log.Printf("SymEncrypt start, key size: %d", me.key.Size())
	return me
}

func (me *SymEncrypter) Write(p []byte) (n int, err error) {
	me.byteCount = me.byteCount + uint64(len(p))
	return me.writer.Write(p)
}

func (me *SymEncrypter) Close() error {
	log.Printf("SymEncrypt done, size: %d", me.byteCount)
	return me.writer.Close()
}

func _symEncryptBytes(plain []byte, key *SymKey) *bytes.Buffer {
	buf := new(bytes.Buffer)
	encrypter := SymEncypterCreate(buf, key)
	w := encrypter.Encrypt()
	defer w.Close()
	w.Write(plain)
	return buf
}

func SymEncryptBytes(plain []byte, key *SymKey) string {
	return _symEncryptBytes(plain, key).String()
}

func SymEncryptString(plain string, key *SymKey) string {
	return SymEncryptBytes([]byte(plain), key)
}
