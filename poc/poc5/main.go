package main

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/openpgp"
)

func main() {
	const publicKey = `
----BEGIN PGP PUBLIC KEY BLOCK-----

xk0EZhMdugECAOXnYV4Kxangww+e/berYrvmrwf38xuQ0OD5FNtKYpSuaq6bqNfn
fhex9RZ4cuA4hXYK77yIqT2at+Ni8ss9nC0AEQEAAQ==
=jXv/
-----END PGP PUBLIC KEY BLOCK-----	
`

	// Convert the public key string to a reader
	reader := strings.NewReader(publicKey)

	// Read the armored public key
	entities, err := openpgp.ReadArmoredKeyRing(reader)
	if err != nil {
		fmt.Printf("Error reading armored public key: %v\n", err)
		return
	}

	// Assuming there is at least one entity in the keyring
	if len(entities) > 0 {
		entity := entities[0]
		fmt.Println("Successfully read PGP public key:")
		for _, identity := range entity.Identities {
			fmt.Printf("Identity: %s\n", identity.Name)
		}
	} else {
		fmt.Println("No entities found in the PGP public key.")
	}
}
