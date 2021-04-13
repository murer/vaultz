package pgp

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockPoc(t *testing.T) {

	ids := []string{"a", "b", "c", "d", "e"}
	np := 3
	positions := make([]int, np)
	for i := 0; i < np; i++ {
		positions[i] = i
	}
	log.Printf("%v", positions)

	current := np - 1

	func() {
		for {
			for positions[current] >= len(ids)+current-np {
				current = current - 1
				if current < 0 {
					return
				}
				if positions[current] < len(ids)+current-np {
					positions[current] = positions[current] + 1

				}
			}
			positions[current] = positions[current] + 1
			log.Printf("%v", positions)
		}
	}()

	assert.Fail(t, "a")
}
