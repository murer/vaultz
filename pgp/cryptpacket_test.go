package pgp

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Printf("RSABits is set to 1024 in test")
	Config.RSABits = 1024
}

func TestPacketDesc(t *testing.T) {
	maria := KeyGenerate("maria", "maria@sample.com")
	bob := KeyGenerate("bob", "bob@sample.com")
	john := KeyGenerate("john", "john@sample.com")

	log.Printf("Key bob: %X", bob.Id())
	log.Printf("Key john: %X", john.Id())

	recipients := KeyRingCreate(bob, john)
	ciphered := EncryptString("mymsg", maria, recipients)
	assert.Equal(t, "mymsg", DecryptString(ciphered, KeyRingCreate(maria), recipients))

	// pkts := packet.NewReader(bytes.NewBuffer(ciphered))
	// for {
	// 	pkt, err := pkts.Next()
	// 	util.Check(err)
	// 	switch p := pkt.(type) {
	// 	case *packet.EncryptedKey:
	// 		log.Printf("EncryptedKey: %#v\n\n", p)
	// 		log.Printf("xxx: %X", p.KeyId)
	// 		abc := recipients.toPgpEntityList().KeysById(p.KeyId)
	// 		for i, kkk := range abc {
	// 			log.Printf("KKK: %d = '%#v'", i, kkk)
	// 		}
	// 	default:
	// 		log.Printf("aaaa: %#v\n\n", pkt)
	// 	}
	// }

}
