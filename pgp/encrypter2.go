package pgp

import "io"

type Encrypter2 struct {
	originalWriter io.Writer
	armored        bool
	recipients     *KeyRing
	signer         *KeyPair
	symKey         *SymKey
}

func CreateEncrypter(writer io.Writer) *Encrypter2 {
	return &Encrypter2{originalWriter: writer}
}

func (me *Encrypter2) Armor(armored bool) *Encrypter2 {
	me.armored = armored
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

func (me *Encrypter2) Start() io.Reader {
	return nil
}
