package pgp

import "golang.org/x/crypto/openpgp"

type KeyRing struct {
	kps map[string]*KeyPair
}

func KeyRingCreate(kps ...*KeyPair) *KeyRing {
	ret := &KeyRing{
		kps: make(map[string]*KeyPair),
	}
	ret.Add(kps...)
	return ret
}

func (me *KeyRing) Get(name string) *KeyPair {
	return me.kps[name]
}

func (me *KeyRing) _add(kp *KeyPair) {
	id := kp.Id()
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
