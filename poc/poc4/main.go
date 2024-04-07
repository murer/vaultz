package main

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/openpgp/armor"
)

func main() {
	buf := new(bytes.Buffer)
	x, err := armor.Encode(buf, "AAAA", nil)
	if err != nil {
		panic(err)
	}
	x.Write([]byte("bla"))
	x.Close()
	fmt.Println(buf.String())
}
