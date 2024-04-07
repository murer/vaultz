package main

import (
	"bytes"
	"crypto"
	"fmt"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var Config = &packet.Config{
	DefaultHash:            crypto.SHA256,
	DefaultCipher:          packet.CipherAES256,
	DefaultCompressionAlgo: packet.CompressionZLIB,
	RSABits:                2048,
	Time: func() time.Time {
		return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	},
}

func main() {
	// Generate a new OpenPGP key pair
	entity, err := openpgp.NewEntity("Example User", "Test Only", "example@example.com", Config)
	if err != nil {
		panic(err)
	}

	// Serialize the public key into a bytes.Buffer
	pubKeyBuf := bytes.NewBuffer(nil)
	err = entity.Serialize(pubKeyBuf)
	if err != nil {
		panic(err)
	}

	// For demonstration, we'll also serialize the public key in armored format
	armoredBuf := bytes.NewBuffer(nil)
	armoredWriter, err := armor.Encode(armoredBuf, openpgp.PublicKeyType, nil)
	if err != nil {
		panic(err)
	}
	err = entity.Serialize(armoredWriter)
	if err != nil {
		panic(err)
	}
	err = armoredWriter.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("Armored Public Key:")
	fmt.Println(armoredBuf.String())

	// Parse the serialized public key back into an entity
	pubKeyReader := bytes.NewReader(pubKeyBuf.Bytes())
	pubEntity, err := openpgp.ReadKeyRing(pubKeyReader)
	if err != nil {
		panic(err)
	}

	// Assuming the first key in the ring is the one we're interested in
	if len(pubEntity) > 0 {
		fmt.Println("Successfully parsed public key back into an entity.")
		fmt.Printf("Entity details: %s <%s>\n", pubEntity[0].PrimaryKey.Fingerprint, pubEntity[0].PrimaryKey.KeyIdString())
	} else {
		fmt.Println("No keys found in the parsed key ring.")
	}
}
