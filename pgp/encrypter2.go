package pgp

import (
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type Encrypter2 struct {
	originalWriter io.Writer
	armored        string
	recipients     *KeyRing
	signer         *KeyPair
	symKey         *SymKey

	writer        io.Writer
	armoredWriter io.WriteCloser
	symWriter     io.WriteCloser
	encryptWriter io.WriteCloser
}

func CreateEncrypter(writer io.Writer) *Encrypter2 {
	return &Encrypter2{originalWriter: writer}
}

func (me *Encrypter2) Armored(blockType string) *Encrypter2 {
	me.armored = blockType
	return me
}

func (me *Encrypter2) Encrypt(recipients *KeyRing) *Encrypter2 {
	me.recipients = recipients
	return me
}

func (me *Encrypter2) Sign(signer *KeyPair) *Encrypter2 {
	me.signer = signer
	return me
}

func (me *Encrypter2) Symmetric(key *SymKey) *Encrypter2 {
	me.symKey = key
	return me
}

func (me *Encrypter2) Start() io.Writer {
	util.Assert(me.originalWriter == nil, "Writer is required")
	util.Assert(me.symKey == nil && me.recipients == nil && me.signer == nil && me.armored == "", "Nothing to do")
	util.Assert(me.symKey != nil && me.recipients != nil, "Symmetric encryption can not have recipients")
	util.Assert(me.symKey != nil && me.signer != nil, "Symmetric encryption can not have signer")

	me.writer = me.originalWriter
	me.preapreArmored()

	if me.symKey == nil && me.recipients == nil && me.signer == nil {
		log.Printf("Encrypter, armor parsing only")
		return &encrypter2Writer{encrypter: me}
	}
	if me.symKey != nil {
		return me.openSymEncrypt()
	}

	return me.openEncrypt()
}

func (me *Encrypter2) getSignerKey() *openpgp.Entity {
	if me.signer == nil {
		return nil
	}
	return me.signer.pgpkey
}

func (me *Encrypter2) openEncrypt() io.Writer {
	encryptWriter, err := openpgp.Encrypt(me.writer, me.recipients.toPgpEntityList(), me.getSignerKey(), nil, nil)
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

func (me *Encrypter2) openSymEncrypt() io.Writer {
	packetConfig := &packet.Config{
		DefaultCipher: packet.CipherAES256,
	}
	symWriter, err := openpgp.SymmetricallyEncrypt(me.writer, me.symKey.key, nil, packetConfig)
	util.Check(err)
	me.symWriter = symWriter
	log.Printf("Encrypt, symmetric with key size: %d", me.symKey.Size())
	me.writer = symWriter
	return &encrypter2Writer{encrypter: me}
}

func (me *Encrypter2) preapreArmored() {
	if me.armored == "" {
		return
	}
	log.Printf("Encrypter, Preparing armor")
	armoredWriter, err := armor.Encode(me.writer, me.armored, nil)
	util.Check(err)
	me.armoredWriter = armoredWriter
	me.writer = armoredWriter
}

type encrypter2Writer struct {
	encrypter *Encrypter2
}

func (me *encrypter2Writer) Write(p []byte) (n int, err error) {
	return me.encrypter.writer.Write(p)
}

func (me *Encrypter2) Close() error {
	if me.encryptWriter != nil {
		log.Printf("Encrypter, closing encrypt writer writer")
		me.encryptWriter.Close()
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
