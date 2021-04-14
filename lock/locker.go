package lock

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/murer/vaultz/pgp"
	"github.com/murer/vaultz/util/combiner"
)

type Locker struct {
	io.WriteCloser
	ciphered   io.Writer
	signer     *pgp.KeyPair
	recipients *pgp.KeyRing
	lockSize   int
	byteCount  uint64
	symKey     *pgp.SymKey
	combiner   *combiner.Combiner
	signWriter io.WriteCloser
	writer     io.WriteCloser
}

func LockerCreate(writer io.Writer, signer *pgp.KeyPair, recipients *pgp.KeyRing, lockSize int) *Locker {
	return &Locker{
		ciphered:   writer,
		signer:     signer,
		recipients: recipients,
		lockSize:   lockSize,
		byteCount:  uint64(0),
		symKey:     pgp.SymKeyGenerate(),
	}
}

func (me *Locker) Write(p []byte) (n int, err error) {
	me.byteCount = me.byteCount + uint64(len(p))
	return me.writer.Write(p)
}

func (me *Locker) Close() error {
	log.Printf("Encrypt done, size: %d", me.byteCount)
	return me.writer.Close()
}

func (me *Locker) writeLocks() {
	ids := me.recipients.Ids()
	me.combiner = combiner.Combine(ids, me.lockSize)
	signer := pgp.SignerCreate(me.ciphered, me.signer)
	me.signWriter = signer.Sign()
	defer me.signWriter.Close()
	header := map[string]string{
		"ls": fmt.Sprintf("%d", me.lockSize),
		"rt": fmt.Sprintf("%d", me.recipients.Size()),
		"cs": fmt.Sprintf("%d", me.combiner.Total()),
	}
	me.signWriter.Write([]byte(stringify(header)))
}

func (me *Locker) Lock() io.WriteCloser {
	me.writeLocks()
	symEncrypter := pgp.SymEncypterCreate(me.ciphered, me.symKey)
	me.writer = symEncrypter.Encrypt()
	return me
}

func LockBytes(data []byte, signer *pgp.KeyPair, recipients *pgp.KeyRing, lockSize int) []byte {
	buf := new(bytes.Buffer)
	func() {
		locker := LockerCreate(buf, signer, recipients, lockSize)
		w := locker.Lock()
		defer locker.Close()
		w.Write(data)
	}()
	return buf.Bytes()
}

func LockString(data string, signer *pgp.KeyPair, recipients *pgp.KeyRing, lockSize int) []byte {
	return LockBytes([]byte(data), signer, recipients, lockSize)
}
