package main

import (
	"bytes"
	"crypto"
	"io"
	"log"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var Config = &packet.Config{
	DefaultHash:            crypto.SHA256,
	DefaultCipher:          packet.CipherAES256,
	DefaultCompressionAlgo: packet.CompressionZLIB,
	CompressionConfig: &packet.CompressionConfig{
		Level: packet.BestCompression,
	},
	RSABits: 512,
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func ArmorIn(writer io.Writer, blockType string) io.WriteCloser {
	ret, err := armor.Encode(writer, openpgp.PublicKeyType, nil)
	Check(err)
	return ret
}

func ArmorInBytes(data []byte) string {
	buf := new(bytes.Buffer)
	func() {
		writer := ArmorIn(buf, openpgp.PublicKeyType)
		defer writer.Close()
		writer.Write(data)
	}()
	return buf.String()
}

func ArmorInPublicKey(key *packet.PublicKey) string {
	buf := new(bytes.Buffer)
	func() {
		writer := ArmorIn(buf, openpgp.PublicKeyType)
		defer writer.Close()
		key.Serialize(writer)
	}()
	return buf.String()
}

func main() {
	fromKP, err := openpgp.NewEntity("John1", "Testing", "johndoe@example.com", Config)
	Check(err)
	log.Printf("From key pair: %v", fromKP.PrimaryKey.KeyIdString())
	log.Println(ArmorInPublicKey(fromKP.PrimaryKey))

	dstKP1, err := openpgp.NewEntity("John2", "Testing", "johndoe@example.com", Config)
	Check(err)
	log.Printf("Dest key pair: %v", dstKP1.PrimaryKey.KeyIdString())
	log.Println(ArmorInPublicKey(dstKP1.PrimaryKey))

	dstKP2, err := openpgp.NewEntity("John3", "Testing", "johndoe@example.com", Config)
	Check(err)
	log.Printf("Dest key pair: %v", dstKP2.PrimaryKey.KeyIdString())
	log.Println(ArmorInPublicKey(dstKP2.PrimaryKey))

	var dests openpgp.EntityList
	dests = append(dests, dstKP1)
	dests = append(dests, dstKP2)

	buf := new(bytes.Buffer)
	log.Printf("Encrypting")
	func() {
		cwriter, err := openpgp.Encrypt(buf, dests, fromKP, nil, Config)
		Check(err)
		defer cwriter.Close()
		cwriter.Write([]byte("mymsg"))
	}()

	log.Printf("Encrypted: %s", ArmorInBytes(buf.Bytes()))

	var keys openpgp.EntityList
	keys = append(keys, dstKP1)
	keys = append(keys, dstKP2)
	keys = append(keys, fromKP)

	log.Printf("Decrypting")
	reader := bytes.NewReader(buf.Bytes())
	msg, err := openpgp.ReadMessage(reader, keys, nil, Config)
	Check(err)
	log.Printf("isSigned: %v", msg.IsSigned)
	log.Printf("SignatureError: %#v", msg.SignatureError)
	log.Printf("SignedByKeyId: %x", msg.SignedByKeyId)
	log.Printf("SignedBy.KeyIdString: %s", msg.SignedBy.PublicKey.KeyIdString())
	log.Printf("SignedBy.Fingerprint: %x", msg.SignedBy.PublicKey.Fingerprint)
	log.Printf("SignedBy.PublicKey: %s", ArmorInPublicKey(msg.SignedBy.PublicKey))
	// log.Printf("SignedByKeyId: %v", msg.SignedBy.PrivateKey)
	log.Printf("IsSymmetricallyEncrypted: %b", msg.IsSymmetricallyEncrypted)
	log.Printf("Signature: %#v", msg.Signature)
	log.Printf("isEncrypted: %v", msg.IsEncrypted)
	buf = &bytes.Buffer{}
	buf.ReadFrom(msg.LiteralData.Body)
	data := buf.String()
	log.Printf("Decrypted: %s", data)

	buf = new(bytes.Buffer)
	log.Printf("Signing")
	func() {
		cwriter, err := openpgp.Sign(buf, fromKP, nil, Config)
		Check(err)
		defer cwriter.Close()
		cwriter.Write([]byte("mymsg"))
	}()

	log.Printf("Signed: %x", buf.Bytes())

	log.Printf("Decrypting")
	reader = bytes.NewReader(buf.Bytes())
	msg, err = openpgp.ReadMessage(reader, keys, nil, Config)
	Check(err)
	log.Printf("isSigned: %v", msg.IsSigned)
	log.Printf("SignatureError: %#v", msg.SignatureError)
	log.Printf("SignedByKeyId: %x", msg.SignedByKeyId)
	log.Printf("SignedBy.KeyIdString: %s", msg.SignedBy.PublicKey.KeyIdString())
	log.Printf("SignedBy.Fingerprint: %x", msg.SignedBy.PublicKey.Fingerprint)
	log.Printf("SignedBy.PublicKey: %s", ArmorInPublicKey(msg.SignedBy.PublicKey))
	// log.Printf("SignedByKeyId: %v", msg.SignedBy.PrivateKey)
	log.Printf("IsSymmetricallyEncrypted: %b", msg.IsSymmetricallyEncrypted)
	log.Printf("Signature: %#v", msg.Signature)
	log.Printf("isEncrypted: %v", msg.IsEncrypted)
	buf = &bytes.Buffer{}
	buf.ReadFrom(msg.LiteralData.Body)
	data = buf.String()
	log.Printf("Decrypted: %s", data)

}
