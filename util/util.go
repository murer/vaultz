package util

import (
	"io"
	"io/ioutil"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadAll(r io.Reader) []byte {
	ret, err := ioutil.ReadAll(r)
	Check(err)
	return ret
}
