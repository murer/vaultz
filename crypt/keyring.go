package crypt

import "golang.org/x/crypto/openpgp"

type KeyRing struct {
	kps []*KeyPair
}

func KeyRingCreate(kps ...*KeyPair) *KeyRing {
	return &KeyRing{kps: kps}
}

func (me *KeyRing) Add(kps ...*KeyPair) *KeyRing {
	me.kps = append(me.kps, kps...)
	return me
}

func (me *KeyRing) toPgpEntityList() openpgp.EntityList {
	var ret openpgp.EntityList
	for _, v := range me.kps {
		ret = append(ret, v.pgpkey)
	}
	return ret
}

func (me *KeyRing) first() *KeyPair {
	return me.kps[0]
}
