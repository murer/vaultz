package main

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func main() {
	key := `
-----BEGIN PGP PUBLIC KEY BLOCK-----

dGVzdA==
=+G7Q
-----END PGP PUBLIC KEY BLOCK-----`
	buf := bytes.NewBuffer([]byte(key))
	x, err := armor.Decode(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(x.Type)
	y, err := openpgp.ReadKeyRing(x.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x: %#v", y)
}
