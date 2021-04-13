package combiner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockPoc(t *testing.T) {

	max := 4
	ids := make([]string, max)
	for i := 0; i < max; i++ {
		ids[i] = fmt.Sprintf("%d", i)
	}
	comb := Combine(ids, 3)
	assert.Equal(t, uint64(4), comb.Total())

	assert.Equal(t, []string{"0", "1", "2"}, comb.Next())
	assert.Equal(t, []string{"0", "1", "3"}, comb.Next())
	assert.Equal(t, []string{"0", "2", "3"}, comb.Next())
	assert.Equal(t, []string{"1", "2", "3"}, comb.Next())
	assert.Equal(t, []string{}, comb.Next())
	assert.Equal(t, []string{}, comb.Next())
}
