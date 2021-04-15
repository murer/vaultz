package util

import (
	"io"
	"io/ioutil"
	"log"
)

func Check(err error) {
	if err != nil {
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
