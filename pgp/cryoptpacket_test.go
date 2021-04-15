package pgp

import "log"

func init() {
	log.Printf("RSABits is set to 1024 in test")
	Config.RSABits = 1024
}
