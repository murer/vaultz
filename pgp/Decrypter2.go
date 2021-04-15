package pgp

import (
	"io"
	"log"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp/armor"
)

type Dectypter2 struct {
	originalReader io.Reader
	armored        bool
	recipients     *KeyRing
	writers        *KeyRing
	SymKey         *SymKey

	armorBlock *armor.Block
	reader     io.Reader
}

type Decryptor2Reader struct {
	decrypter *Dectypter2
}

func (me *Decryptor2Reader) Read(p []byte) (n int, err error) {
	return me.decrypter.reader.Read(p)
}

func (me *Decryptor2Reader) Close() error {
	return nil
}

func CreateDecrypter(reader io.Reader) *Dectypter2 {
	return &Dectypter2{originalReader: reader}
}

func (me *Dectypter2) Armor(armored bool) *Dectypter2 {
	me.armored = armored
	return me
}

func (me *Dectypter2) Decrypt(recipients *KeyRing) *Dectypter2 {
	me.recipients = recipients
	return me
}

func (me *Dectypter2) Verify(writers *KeyRing) *Dectypter2 {
	me.writers = writers
	return me
}

func (me *Dectypter2) Symmetric(key *SymKey) *Dectypter2 {
	me.SymKey = key
	return me
}

func (me *Dectypter2) Open() io.ReadCloser {
	me.check(me.originalReader == nil, "Reader is required")
	me.check(me.SymKey == nil && me.recipients == nil && me.writers == nil && !me.armored, "Nothing to do")
	me.check(me.SymKey != nil && me.recipients != nil, "Symmetric decryption can not have recipients")
	me.check(me.SymKey != nil && me.writers != nil, "Symmetric decryption can not have writers")

	me.reader = me.originalReader
	me.preapreArmored()
	if me.SymKey == nil && me.recipients == nil && me.writers == nil {
		return me.openArmored()
	}

	if me.SymKey != nil {
		return me.openSymDecrypt()
	}

	return nil
}

func (me *Dectypter2) openArmored() io.ReadCloser {
	log.Printf("Decrypter, armor parsing only")
	return &Decryptor2Reader{decrypter: me}
}

func (me *Dectypter2) openSymDecrypt() io.ReadCloser {
	return nil
}

func (me *Dectypter2) check(cond bool, msg string) {
	if cond {
		log.Panicf(msg)
	}
}

func (me *Dectypter2) preapreArmored() {
	if !me.armored {
		return
	}
	log.Printf("Decrypter, Prepareing armor parsing")
	block, err := armor.Decode(me.originalReader)
	util.Check(err)
	me.armorBlock = block
	me.reader = block.Body
}
