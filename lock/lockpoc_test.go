package pgp

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type combinator struct {
	ids        []string
	tuple      []int
	hasStarted bool
}

func (me *combinator) checkPosition(position int, value int) {
	max := len(me.ids) + position - len(me.tuple)
	if value > max || value < 0 {
		log.Panicf("wrong for position: %d, max: %d, but it is: %d, len ids: %d, len tuple: %d", position, max, value, len(me.ids), len(me.tuple))
	}
}

func (me *combinator) setPosition(position int, value int) {
	me.hasStarted = true
	for i := position; i < len(me.tuple); i++ {
		nv := value + i - position
		me.checkPosition(i, nv)
		me.tuple[i] = nv
	}
}

func (me *combinator) detuple() []string {
	ret := make([]string, len(me.tuple))
	for i := 0; i < len(me.tuple); i++ {
		ret[i] = me.ids[me.tuple[i]]
	}
	return ret
}

func (me *combinator) Next() []string {
	if !me.hasStarted {
		me.setPosition(0, 0)
		return me.detuple()
	}
	for i := len(me.tuple) - 1; i >= 0; i-- {
		if me.tuple[i] < len(me.ids)+i-len(me.tuple) {
			me.setPosition(i, me.tuple[i]+1)
			return me.detuple()
		}
	}
	return []string{}
}

func TestLockPoc(t *testing.T) {

	comb := &combinator{
		ids:        []string{"a", "b", "c", "d", "e"},
		tuple:      make([]int, 3),
		hasStarted: false,
	}
	for i := 0; i < 10; i++ {
		fmt.Printf("x: %d = %v\n", i, comb.Next())
	}

	assert.Fail(t, "a")
}
