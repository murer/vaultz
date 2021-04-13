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

func (me *combinator) Total() uint64 {
	x := len(me.ids) - len(me.tuple)
	ret := uint64(1)
	for i := x + 1; i <= len(me.ids); i++ {
		ret = ret * uint64(i)
	}
	y := uint64(1)
	for i := 1; i <= len(me.tuple); i++ {
		y = y * uint64(i)
	}
	return ret / y
}

func TestLockPoc(t *testing.T) {

	comb := &combinator{
		tuple:      make([]int, 2),
		hasStarted: false,
	}
	max := 40
	comb.ids = make([]string, max)
	for i := 0; i < max; i++ {
		comb.ids[i] = fmt.Sprintf("%d", i)
	}
	fmt.Printf("total: %d\n", comb.Total())
	for i := comb.Next(); len(i) > 0; i = comb.Next() {
		fmt.Printf("x: %v\n", i)
	}

	assert.Fail(t, "a")
}
