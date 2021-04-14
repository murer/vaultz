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
	ciphered := LockString("mymsg", ring.Get("u"), ring, 3)
	log.Printf("Ciphered\n%s", pgp.ArmorEncodeBytes(ciphered, "PGP MESSAGE"))

}
