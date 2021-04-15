package lock

// func TestPocLock(t *testing.T) {
// 	signer := pgp.KeyGenerate("u0", "u0@sample.com")
// 	ring := pgp.KeyRingCreate(signer)
// 	for i := 1; i < 5; i++ {
// 		ring.Add(pgp.KeyGenerate(fmt.Sprintf("u%d", i), fmt.Sprintf("u%d@sample.com", i)))
// 	}
// 	ciphered := LockString("mymsg", signer, ring, 3)
// 	log.Printf("Ciphered\n%s", pgp.ArmorEncodeBytes(ciphered, "PGP MESSAGE"))

// }
