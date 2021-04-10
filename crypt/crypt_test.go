package crypt

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	"github.com/stretchr/testify/assert"
)

func TestPoc(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := &KeyPair{}
	kp.Generate("test", "test@sample.com")

	// Add more identities here if you wish

	// Sign all the identities
	// for _, id := range e.Identities {
	// 	err := id.SelfSignature.SignUserId(id.UserId.Id, e.PrimaryKey, e.PrivateKey, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }

	// w, err := armor.Encode(os.Stdout, openpgp.PublicKeyType, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer w.Close()

	// e.Serialize(w)

	s, err := armor.Encode(os.Stdout, openpgp.PrivateKeyType, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()

	kp.pgpkey.SerializePrivate(s, nil)
}
