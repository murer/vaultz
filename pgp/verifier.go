package pgp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func VerifierCreate(plain io.Reader, writers *KeyRing) *Verifier {
	return &Verifier{plain: plain, writers: writers}
}

type Verifier struct {
	io.Closer
	plain   io.Reader
	writers *KeyRing

	msg      *openpgp.MessageDetails
	tempKey  *SymKey
	tempFile string
}

func (me *Verifier) Close() error {
	if me.tempFile != "" {
		log.Printf("Verifier closing, delete: %s", me.tempFile)
		os.Remove(me.tempFile)
	}
	return nil
}

func (me *Verifier) UnsafeDecrypt() io.Reader {
	ar, err := armor.Decode(me.plain)
	util.Check(err)
	msg, err := openpgp.ReadMessage(ar.Body, me.writers.toPgpEntityList(), nil, nil)
	util.Check(err)
	me.msg = msg
	if me.msg.IsEncrypted {
		log.Panicf("Verifier, it is actually encrypted")
	}
	if !me.msg.IsSigned {
		log.Panicf("Verifier, it is not signed")
	}
	return me.msg.UnverifiedBody
}

func (me *Verifier) UnsafeDecryptBytes() []byte {
	r := me.UnsafeDecrypt()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	return ret
}

func (me *Verifier) UnsafeDecryptString() string {
	return string(me.UnsafeDecryptBytes())
}

func (me *Verifier) decryptToTemp() {
	// TODO: must be encrypted
	unsafe := me.UnsafeDecrypt()
	f, err := ioutil.TempFile(os.TempDir(), "vaultz-decrypt-*.tmp")
	util.Check(err)
	defer f.Close()
	me.tempKey = SymKeyGenerate()
	encrypter := SymEncypterCreate(f, me.tempKey)
	defer encrypter.Close()
	total, err := io.Copy(encrypter.Encrypt(), unsafe)
	util.Check(err)
	log.Printf("Verifiy to %s, total: %d", f.Name(), total)
	me.tempFile = f.Name()
}

func (me *Verifier) decryptCheckSigner() {
	if me.msg.Signature == nil {
		msg := fmt.Sprintf("Verify unknown signer: %X", me.msg.SignedByKeyId)
		log.Panic(msg)
	}
	sigKP := keyFromEntity(me.msg.SignedBy.Entity)
	pubKey := sigKP.ExportPub()
	for _, v := range me.writers.kps {
		if v.ExportPub() == pubKey {
			return
		}
	}
	log.Panicf("Verify signer is not a writer: %s", sigKP.Id())
}

func (me *Verifier) Decrypt() io.ReadCloser {
	me.decryptToTemp()
	me.decryptCheckSigner()
	signerKP := keyFromEntity(me.msg.SignedBy.Entity)
	log.Printf("Verify signed by id: %s %s", signerKP.Id(), signerKP.UserName())
	f, err := os.Open(me.tempFile)
	util.Check(err)
	Verifier := SymDecrypterCreate(f, me.tempKey)
	return Verifier.Decrypt()
}

func (me *Verifier) DecryptBytes() []byte {
	r := me.Decrypt()
	defer r.Close()
	ret, err := ioutil.ReadAll(r)
	util.Check(err)
	return ret
}

func (me *Verifier) DecryptString() string {
	return string(me.DecryptBytes())
}
