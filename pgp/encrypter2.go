package pgp

import (
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp/armor"
)

type Encrypter2 struct {
	originalWriter io.Writer
	armored        string
	recipients     *KeyRing
	signer         *KeyPair
	symKey         *SymKey

	writer        io.Writer
	armoredWriter io.WriteCloser
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

	return nil
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
	if me.armoredWriter != nil {
		log.Printf("Encrypter, closing armored writer")
		me.armoredWriter.Close()
	}
	return nil
}
