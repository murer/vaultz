package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoc(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := &KeyPair{}
	kp.Generate("test", "test@sample.com")
	assert.NotNil(t, kp.ExportPub())
	assert.NotNil(t, kp.ExportPriv())

	pubkp := &KeyPair{}
	pubkp.Import(kp.ExportPub())
	assert.Equal(t, kp.ExportPub(), pubkp.ExportPub())
	assert.Equal(t, "", pubkp.ExportPriv())

	privkp := &KeyPair{}
	privkp.Import(kp.ExportPriv())
	assert.Equal(t, kp.ExportPub(), privkp.ExportPub())
	assert.Equal(t, kp.ExportPriv(), privkp.ExportPriv())

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
