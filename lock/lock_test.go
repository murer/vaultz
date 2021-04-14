package lock

import (
	"fmt"
	"log"
	"testing"

	"github.com/murer/vaultz/pgp"
)

func TestPocLock(t *testing.T) {
	ring := pgp.KeyRingCreate()
	for i := 0; i < 5; i++ {
		ring.Add(pgp.KeyGenerate(fmt.Sprintf("u%d", i), fmt.Sprintf("u%d@sample.com", i)))
	}
	locker := &Locker{
		signer:     ring.Get("x"),
		recipients: ring,
		lockSize:   3,
	}
	ciphered := locker.LockString("mymsg")
	log.Printf("Ciphered %X", pgp.ArmorEncodeBytes(ciphered, "VAULTZ CONTENT"))

}
