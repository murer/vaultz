package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexUInt64(t *testing.T) {
	assert.Equal(t, "000000000000000A", HexUInt64(10))
	assert.Equal(t, "000000000000007F", HexUInt64(127))
}
