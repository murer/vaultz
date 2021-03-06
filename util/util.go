package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Checkf(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Printf(msg, args...)
		panic(err)
	}
}

func Assert(cond bool, msg string, v ...interface{}) {
	if cond {
		log.Panicf(msg, v...)
	}
}

func ReadAll(r io.Reader) []byte {
	ret, err := ioutil.ReadAll(r)
	Check(err)
	return ret
}

func ReadAllString(r io.Reader) string {
	return string(ReadAll(r))
}

func HexUInt64(n uint64) string {
	return fmt.Sprintf("%016X", n)
}
