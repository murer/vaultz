package crypt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// fmt.Println(maria.ExportPub())
	// fmt.Println(bob.ExportPriv())

	ring := KeyRingCreate(maria, bob, john)

	ciphered := EncryptString("mymsg", maria, ring)
	// fmt.Println(ciphered)

	ring = KeyRingCreate(maria.PubOnly(), john)
	decrypter := DecrypterCreate(strings.NewReader(ciphered), ring)
	unsafePlain := decrypter.UnsafeDecryptString()
	assert.Equal(t, "mymsg", unsafePlain)

	decrypter = DecrypterCreate(strings.NewReader(ciphered), ring)
	plain := decrypter.DecryptString()
	assert.Equal(t, "mymsg", plain)

	// buf := new(bytes.Buffer)
	// w, err := openpgp.Encrypt(buf, ring.toPgpEntityList(), maria.pgpkey, nil, nil)
	// util.Check(err)
	// w.Write([]byte("mymsg"))
	// w.Close()

	// fbuf, err := os.Create("/tmp/x.txt.pgp")
	// util.Check(err)
	// fbuf.Write(buf.Bytes())
	// fbuf.Close()

	// m, err := openpgp.ReadMessage(buf, ring.toPgpEntityList(), nil, nil)
	// util.Check(err)
	// fmt.Printf("x: %#v\n", m)

	// fmt.Printf("IsEncrypted: %t\n", m.IsEncrypted)
	// fmt.Printf("IsSigned: %t\n", m.IsSigned)
	// fmt.Printf("IsSymmetricallyEncrypted: %t\n", m.IsSymmetricallyEncrypted)

	// data, err := ioutil.ReadAll(m.UnverifiedBody)
	// util.Check(err)
	// fmt.Printf("SignedByKeyId: %s\n", string(data))
	// fmt.Printf("SignatureError: %v\n", m.SignatureError)
	// fmt.Printf("Signature: %#v\n", m.Signature)

	// fmt.Printf("SignedBy: %#v\n", m.SignedBy)
	// assert.Equal(t, maria.pgpkey.PrimaryKey.Fingerprint, m.SignedBy.Entity.PrimaryKey.Fingerprint)

}
