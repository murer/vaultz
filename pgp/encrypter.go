package pgp

import (
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type Encrypter struct {
	originalWriter io.Writer
	armored        string
	recipients     *KeyRing
	signer         *KeyPair
	symKey         *SymKey

	writer        io.Writer
	armoredWriter io.WriteCloser
	symWriter     io.WriteCloser
	encryptWriter io.WriteCloser
	signerWriter  io.WriteCloser
}

func CreateEncrypter(writer io.Writer) *Encrypter {
	return &Encrypter{originalWriter: writer}
}

func (me *Encrypter) Armored(blockType string) *Encrypter {
	me.armored = blockType
	return me
}

func (me *Encrypter) Encrypt(recipients *KeyRing) *Encrypter {
	me.recipients = recipients.PubOnly()
	return me
}

func (me *Encrypter) Sign(signer *KeyPair) *Encrypter {
	me.signer = signer
	return me
}

func (me *Encrypter) Symmetric(key *SymKey) *Encrypter {
	me.symKey = key
	return me
}

func (me *Encrypter) Start() io.Writer {
	util.Assert(me.originalWriter == nil, "Writer is required")
	util.Assert(me.symKey == nil && me.recipients == nil && me.signer == nil && me.armored == "", "Nothing to do")
	util.Assert(me.symKey != nil && me.recipients != nil, "Symmetric encryption can not have recipients")
	util.Assert(me.symKey != nil && me.signer != nil, "Symmetric encryption can not have signer")

	me.writer = me.originalWriter
	me.preapreArmored()

	if me.symKey == nil && me.recipients == nil && me.signer == nil {
		log.Printf("Encrypter, armor parsing only")
		return &encrypterWriter{encrypter: me}
	}
	if me.symKey != nil {
		return me.openSymEncrypt()
	}
	if me.recipients == nil {
		return me.openSigner()
	}

	return me.openEncrypt()
}

func (me *Encrypter) getSignerKey() *openpgp.Entity {
	if me.signer == nil {
		return nil
	}
	return me.signer.pgpkey
}

func (me *Encrypter) openSigner() io.Writer {
	signerWriter, err := openpgp.Sign(me.writer, me.getSignerKey(), nil, Config)
	util.Check(err)
	me.signerWriter = signerWriter
	if me.getSignerKey() != nil {
		log.Printf("Encrypt, signer only: %s %s", me.signer.Id(), me.signer.UserName())
	}
	me.writer = signerWriter
	return signerWriter
}

func (me *Encrypter) openEncrypt() io.Writer {
	encryptWriter, err := openpgp.Encrypt(me.writer, me.recipients.toPgpEntityList(), me.getSignerKey(), nil, Config)
	util.Check(err)
	me.encryptWriter = encryptWriter
	if me.getSignerKey() != nil {
		log.Printf("Encrypt, signer: %s %s, total recipients: %d", me.signer.Id(), me.signer.UserName(), len(me.recipients.kps))
	} else {
		log.Printf("Encrypt, no signer, total recipients: %d", len(me.recipients.kps))
	}
	for _, v := range me.recipients.kps {
		log.Printf("Encrypt, recipients: %s %s", v.Id(), v.UserName())
	}
	me.writer = encryptWriter
	return encryptWriter
}

func (me *Encrypter) openSymEncrypt() io.Writer {
	packetConfig := &packet.Config{
		DefaultCipher: packet.CipherAES256,
	}
	symWriter, err := openpgp.SymmetricallyEncrypt(me.writer, me.symKey.key, nil, packetConfig)
	util.Check(err)
	me.symWriter = symWriter
	log.Printf("Encrypt, symmetric with key size: %d", me.symKey.Size())
	me.writer = symWriter
	return &encrypterWriter{encrypter: me}
}

func (me *Encrypter) preapreArmored() {
	if me.armored == "" {
		return
	}
	log.Printf("Encrypter, Preparing armor")
	armoredWriter, err := armor.Encode(me.writer, me.armored, nil)
	util.Check(err)
	me.armoredWriter = armoredWriter
	me.writer = armoredWriter
}

type encrypterWriter struct {
	encrypter *Encrypter
}

func (me *encrypterWriter) Write(p []byte) (n int, err error) {
	return me.encrypter.writer.Write(p)
}

func (me *Encrypter) Close() error {
	if me.encryptWriter != nil {
		log.Printf("Encrypter, closing encrypt writer writer")
		me.encryptWriter.Close()
	}
	if me.signerWriter != nil {
		log.Printf("Encrypter, closing encrypt writer writer")
		me.signerWriter.Close()
	}
	if me.symWriter != nil {
		log.Printf("Encrypter, closing symmetric writer writer")
		me.symWriter.Close()
	}
	if me.armoredWriter != nil {
		log.Printf("Encrypter, closing armored writer")
		me.armoredWriter.Close()
	}
	return nil
}
