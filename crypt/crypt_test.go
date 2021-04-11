package crypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGen(t *testing.T) {
	assert.Equal(t, 1, 1)

	kp := &KeyPair{}
	kp.Generate("test", "test@sample.com")
	assert.NotEmpty(t, kp.ExportPub())
	assert.NotEmpty(t, kp.ExportPriv())
	assert.NotEmpty(t, kp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	pubkp := &KeyPair{}
	pubkp.Import(kp.ExportPub())
	assert.Equal(t, kp.ExportPub(), pubkp.ExportPub())
	assert.Empty(t, pubkp.ExportPriv())
	assert.Equal(t, kp.Id(), pubkp.Id())
	fmt.Printf("id: %s\n", kp.Id())

	privkp := &KeyPair{}
	privkp.Import(kp.ExportPriv())
	assert.Equal(t, kp.ExportPub(), privkp.ExportPub())
	assert.NotEmpty(t, privkp.ExportPriv())
	assert.Equal(t, kp.Id(), privkp.Id())
	fmt.Printf("id: %s\n", privkp.Id())
}
