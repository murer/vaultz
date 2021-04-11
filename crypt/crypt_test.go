package crypt

import (
	"bytes"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"

	"fmt"
	"testing"

	"github.com/murer/vaultz/crypt/util"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/openpgp"
)

func TestKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := KeyGenerate("test", "test@sample.com")
	assert.Equal(t, "test", kp.UserName())
	assert.Equal(t, "test@sample.com", kp.UserEmail())
	assert.NotEmpty(t, kp.ExportPub())
	assert.NotEmpty(t, kp.ExportPriv())
	assert.NotEmpty(t, kp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	pubkp := KeyImport(kp.ExportPub())
	assert.Equal(t, "test", pubkp.UserName())
	assert.Equal(t, "test@sample.com", pubkp.UserEmail())
	assert.Equal(t, kp.ExportPub(), pubkp.ExportPub())
	assert.Empty(t, pubkp.ExportPriv())
	assert.Equal(t, kp.Id(), pubkp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	privkp := KeyImport(kp.ExportPriv())
	assert.Equal(t, "test", privkp.UserName())
	assert.Equal(t, "test@sample.com", privkp.UserEmail())
	assert.Equal(t, kp.ExportPub(), privkp.ExportPub())
	assert.NotEmpty(t, privkp.ExportPriv())
	assert.Equal(t, kp.Id(), privkp.Id())
	fmt.Printf("id: %s\n", privkp.Id())
}

func TestCrypt(t *testing.T) {

	maria := KeyGenerate("maria", "maria@sample.com")
	fmt.Printf("maria: %s\n", maria.Id())
	bob := KeyGenerate("bob", "bob@sample.com")
	fmt.Printf("bob: %s\n", bob.Id())
	john := KeyGenerate("john", "john@sample.com")
	fmt.Printf("john: %s\n", john.Id())

	buf := new(bytes.Buffer)
	// buf, err := os.Create("/tmp/x.txt.pgp")
	// util.Check(err)
	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{bob.pgpkey, john.pgpkey}, maria.pgpkey, nil, nil)
	util.Check(err)
	w.Write([]byte("mymsg"))
	w.Close()

	buf2 := new(bytes.Buffer)
	a, err := armor.Encode(buf2, openpgp.PrivateKeyType, nil)
	util.Check(err)
	john.pgpkey.SerializePrivate(a, nil)
	maria.pgpkey.Serialize(a)
	a.Close()
	krs := buf2.String()

	ring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(krs))
	util.Check(err)
	m, err := openpgp.ReadMessage(buf, ring, nil, nil)
	util.Check(err)
	fmt.Printf("x: %#v\n", m)

	fmt.Printf("IsEncrypted: %t\n", m.IsEncrypted)
	fmt.Printf("IsSigned: %t\n", m.IsSigned)
	fmt.Printf("IsSymmetricallyEncrypted: %t\n", m.IsSymmetricallyEncrypted)

	data, err := ioutil.ReadAll(m.UnverifiedBody)
	util.Check(err)
	fmt.Printf("SignedByKeyId: %s\n", string(data))
	fmt.Printf("SignatureError: %v\n", m.SignatureError)
	fmt.Printf("Signature: %v\n", m.Signature)

}
