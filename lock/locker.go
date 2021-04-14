package lock

import "github.com/murer/vaultz/pgp"

type Locker struct {
	signer     *pgp.KeyPair
	recipients *pgp.KeyRing
	lockSize   int
}

func (me *Locker) LockString(data string) []byte {
	return []byte{1}
}
