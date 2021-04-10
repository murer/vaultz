package crypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoc(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := &KeyPair{}
	kp.Generate("test", "test@sample.com")
	fmt.Printf("pub_: %s\n", kp.PubFingerprint())
	fmt.Printf("priv: %s\n", kp.PrivFingerprint())
	fmt.Printf("id__: %s\n", kp.pgpkey.PrivateKey.KeyIdString())
	assert.NotNil(t, kp.ExportPub())
	assert.NotNil(t, kp.ExportPriv())

	pubkp := &KeyPair{}
	pubkp.Import(kp.ExportPub())
	fmt.Printf("pub_: %s\n", pubkp.PubFingerprint())
	fmt.Printf("id__: %s\n", pubkp.pgpkey.PrimaryKey.KeyIdString())
	// fmt.Printf("priv: %s\n", pubkp.PrivFingerprint())
	// assert.Equal(t, kp.ExportPub(), pubkp.ExportPub())
	// assert.Equal(t, "", pubkp.ExportPriv())

	privkp := &KeyPair{}
	privkp.Import(kp.ExportPriv())
	fmt.Printf("pub_: %s\n", privkp.PubFingerprint())
	fmt.Printf("priv: %s\n", privkp.PrivFingerprint())
	fmt.Printf("id__: %s\n", privkp.pgpkey.PrivateKey.KeyIdString())
	// assert.Equal(t, kp.ExportPub(), privkp.ExportPub())
	// assert.Equal(t, kp.ExportPriv(), privkp.ExportPriv())

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

	// kp.pgpkey.Serialize(w)

	// s, err := armor.Encode(os.Stdout, openpgp.PrivateKeyType, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer s.Close()

	// kp.pgpkey.SerializePrivate(s, nil)
}
