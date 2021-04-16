package lock

// type Locker struct {
// 	originalWriter io.Writer
// 	armored        bool
// 	recipients     *pgp.KeyRing
// 	signer         *pgp.KeyPair

// 	writer        io.Writer
// 	armoredWriter io.WriteCloser
// }

// func CreateLocker(writer io.Writer) *Locker {
// 	return &Locker{originalWriter: writer}
// }

// func (me *Locker) Armored(armored bool) *Locker {
// 	me.armored = armored
// 	return me
// }

// func (me *Locker) Encrypt(recipients *pgp.KeyRing) *Locker {
// 	me.recipients = recipients.PubOnly()
// 	return me
// }

// func (me *Locker) Sign(signer *pgp.KeyPair) *Locker {
// 	me.signer = signer
// 	return me
// }

// func (me *Locker) Start() io.Writer {
// 	util.Assert(me.originalWriter == nil, "Writer is required")
// 	util.Assert(me.recipients == nil, "Recipients is required")
// 	util.Assert(me.signer == nil, "Signer is required")

// 	me.writer = me.originalWriter
// 	me.preapreArmored()

// 	if me.symKey == nil && me.recipients == nil && me.signer == nil {
// 		log.Printf("Locker, armor parsing only")
// 		return &LockerWriter{Locker: me}
// 	}
// 	if me.symKey != nil {
// 		return me.openSymEncrypt()
// 	}
// 	if me.recipients == nil {
// 		return me.openSigner()
// 	}

// 	return me.openEncrypt()
// }

// func (me *Locker) getSignerKey() *openpgp.Entity {
// 	if me.signer == nil {
// 		return nil
// 	}
// 	return me.signer.pgpkey
// }

// func (me *Locker) openSigner() io.Writer {
// 	signerWriter, err := openpgp.Sign(me.writer, me.getSignerKey(), nil, Config)
// 	util.Check(err)
// 	me.signerWriter = signerWriter
// 	if me.getSignerKey() != nil {
// 		log.Printf("Encrypt, signer only: %X %s", me.signer.Id(), me.signer.UserName())
// 	}
// 	me.writer = signerWriter
// 	return signerWriter
// }

// func (me *Locker) openEncrypt() io.Writer {
// 	encryptWriter, err := openpgp.Encrypt(me.writer, me.recipients.toPgpEntityList(), me.getSignerKey(), nil, Config)
// 	util.Check(err)
// 	me.encryptWriter = encryptWriter
// 	if me.getSignerKey() != nil {
// 		log.Printf("Encrypt, signer: %X %s, total recipients: %d", me.signer.Id(), me.signer.UserName(), len(me.recipients.kps))
// 	} else {
// 		log.Printf("Encrypt, no signer, total recipients: %d", len(me.recipients.kps))
// 	}
// 	for _, v := range me.recipients.kps {
// 		log.Printf("Encrypt, recipients: %X %s", v.Id(), v.UserName())
// 	}
// 	me.writer = encryptWriter
// 	return encryptWriter
// }

// func (me *Locker) openSymEncrypt() io.Writer {
// 	symWriter, err := openpgp.SymmetricallyEncrypt(me.writer, me.symKey.key, nil, Config)
// 	util.Check(err)
// 	me.symWriter = symWriter
// 	log.Printf("Encrypt, symmetric with key size: %d", me.symKey.Size())
// 	me.writer = symWriter
// 	return &LockerWriter{Locker: me}
// }

// func (me *Locker) preapreArmored() {
// 	if me.armored == "" {
// 		return
// 	}
// 	log.Printf("Locker, Preparing armor")
// 	armoredWriter, err := armor.Encode(me.writer, me.armored, nil)
// 	util.Check(err)
// 	me.armoredWriter = armoredWriter
// 	me.writer = armoredWriter
// }

// type LockerWriter struct {
// 	Locker *Locker
// }

// func (me *LockerWriter) Write(p []byte) (n int, err error) {
// 	return me.Locker.writer.Write(p)
// }

// func (me *Locker) Close() error {
// 	if me.encryptWriter != nil {
// 		log.Printf("Locker, closing encrypt writer writer")
// 		me.encryptWriter.Close()
// 	}
// 	if me.signerWriter != nil {
// 		log.Printf("Locker, closing encrypt writer writer")
// 		me.signerWriter.Close()
// 	}
// 	if me.symWriter != nil {
// 		log.Printf("Locker, closing symmetric writer writer")
// 		me.symWriter.Close()
// 	}
// 	if me.armoredWriter != nil {
// 		log.Printf("Locker, closing armored writer")
// 		me.armoredWriter.Close()
// 	}
// 	return nil
// }
