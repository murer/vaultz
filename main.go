package main

import (
	"bytes"
	"crypto"
	"log"

	"golang.org/x/crypto/openpgp"
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

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fromKP, err := openpgp.NewEntity("John1", "Testing", "johndoe@example.com", Config)
	Check(err)
	log.Printf("From key pair: %v", fromKP.PrimaryKey.KeyIdString())

	dstKP1, err := openpgp.NewEntity("John2", "Testing", "johndoe@example.com", Config)
	Check(err)
	log.Printf("Dest key pair: %v", dstKP1.PrimaryKey.KeyIdString())

	dstKP2, err := openpgp.NewEntity("John3", "Testing", "johndoe@example.com", Config)
	Check(err)
	log.Printf("Dest key pair: %v", dstKP2.PrimaryKey.KeyIdString())
	var dests openpgp.EntityList

	dests = append(dests, dstKP1)
	dests = append(dests, dstKP2)

	buf := new(bytes.Buffer)
	// func() {
	// 	encrypter := CreateEncrypter(buf).Sign(signer).Recipients(recipients)
	// 	defer encrypter.Close()
	// 	encrypter.Start().Write(plain)
	// }()
	// return buf.Bytes()

	func() {
		cwriter, err := openpgp.Encrypt(buf, dests, fromKP, nil, Config)
		Check(err)
		defer cwriter.Close()
		cwriter.Write([]byte("mymsg"))
	}()

}
