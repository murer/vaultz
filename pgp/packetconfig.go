package pgp

import (
	"crypto"

	"golang.org/x/crypto/openpgp/packet"
)

var Config = &packet.Config{
	DefaultHash:            crypto.SHA256,
	DefaultCipher:          packet.CipherAES256,
	DefaultCompressionAlgo: packet.CompressionZLIB,
	CompressionConfig: &packet.CompressionConfig{
		Level: packet.BestCompression,
	},
	RSABits: 1024,
}
