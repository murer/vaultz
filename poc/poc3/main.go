package main

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/openpgp"
)

func main() {
	key := `
-----BEGIN PGP PUBLIC KEY BLOCK-----

dGVzdA==
=+G7Q
-----END PGP PUBLIC KEY BLOCK-----`
	buf := bytes.NewBuffer([]byte(key))
	// x, err := armor.Decode(buf)
	x, err := openpgp.ReadArmoredKeyRing(buf)
	if err != nil {
		panic(err)
	}
	// fmt.Println(x.Type)
	fmt.Printf("x: %v", x)
}
