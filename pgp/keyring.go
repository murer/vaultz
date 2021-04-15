package pgp

import (
	"log"
	"sort"

	"golang.org/x/crypto/openpgp"
)

type KeyRing struct {
	kps map[uint64]*KeyPair
}

func KeyRingCreate(kps ...*KeyPair) *KeyRing {
	ret := &KeyRing{
		kps: make(map[uint64]*KeyPair),
	}
	ret.Add(kps...)
	return ret
}

func (me *KeyRing) IdString() []string {
	ret := make([]string, me.Size())
	i := 0
	for _, v := range me.kps {
		ret[i] = v.Id()
		i = i + 1
	}
	sort.Strings(ret)
	return ret
}

func (me *KeyRing) Ids() []uint64 {
	ret := make([]uint64, me.Size())
	i := 0
	for _, v := range me.kps {
		ret[i] = v.IdBinary()
		i = i + 1
	}
	sort.Slice(ret, func(a int, b int) bool { return ret[a] < ret[b] })
	return ret
}

func (me *KeyRing) Size() int {
	return len(me.kps)
}

func (me *KeyRing) Get(id uint64) *KeyPair {
	return me.kps[id]
}

func (me *KeyRing) _add(kp *KeyPair) {
	id := kp.IdBinary()
	if me.kps[id] != nil {
		log.Panicf("KeyId collision in the ring: %X", id)
	}
	me.kps[id] = kp
}

func (me *KeyRing) Add(kps ...*KeyPair) *KeyRing {
	for _, v := range kps {
		me._add(v)
	}
	return me
}

func (me *KeyRing) toPgpEntityList() openpgp.EntityList {
	var ret openpgp.EntityList
	for _, v := range me.kps {
		ret = append(ret, v.pgpkey)
	}
	return ret
}

func (me *KeyRing) PubOnly() *KeyRing {
	ret := KeyRingCreate()
	for _, v := range me.kps {
		ret.Add(v.PubOnly())
	}
	return ret
}
