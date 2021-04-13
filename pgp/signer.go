package pgp

import (
	"bytes"
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type Signer struct {
	io.WriteCloser

	ciphered io.Writer
	signer   *KeyPair

	armor  io.WriteCloser
	writer io.WriteCloser

	byteCount uint64
}

func SignerCreate(ciphered io.Writer, signer *KeyPair) *Signer {
	return &Signer{
		ciphered: ciphered,
		signer:   signer,
	}
}

func (me *Signer) Sign() io.WriteCloser {
	wa, err := armor.Encode(me.ciphered, "PGP MESSAGE", nil)
	util.Check(err)
	writer, err := openpgp.Sign(wa, me.signer.pgpkey, nil, nil)
	util.Check(err)
	me.armor = wa
	me.writer = writer
	log.Printf("Signer start, signer: %s %s", me.signer.Id(), me.signer.UserName())
	return me
}

func (me *Signer) Write(p []byte) (n int, err error) {
	me.byteCount = me.byteCount + uint64(len(p))
	return me.writer.Write(p)
}

func (me *Signer) Close() error {
	log.Printf("Signer done, size: %d", me.byteCount)
	we := me.writer.Close()
	ae := me.armor.Close()
	if we != nil {
		return we
	}
	return ae
}

func _signBytes(plain []byte, signer *KeyPair) *bytes.Buffer {
	buf := new(bytes.Buffer)
	sr := SignerCreate(buf, signer)
	w := sr.Sign()
	defer w.Close()
	w.Write(plain)
	return buf
}

func SignBytes(plain []byte, signer *KeyPair) string {
	return _signBytes(plain, signer).String()
}

func SignString(plain string, signer *KeyPair) string {
	return SignBytes([]byte(plain), signer)
}
