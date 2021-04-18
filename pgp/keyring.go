package pgp

import (
	"io"
	"log"
	"sort"

	"github.com/murer/vaultz/util"
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
		ret[i] = v.IdString()
		i = i + 1
	}
	sort.Strings(ret)
	return ret
}

func (me *KeyRing) Ids() []uint64 {
	ret := make([]uint64, me.Size())
	i := 0
	for _, v := range me.kps {
		ret[i] = v.Id()
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
	id := kp.Id()
	old := me.kps[id]
	if old != nil {
		log.Panicf("KeyId: %016X collision in the ring, old: %s, new: %s", id, old.UserName(), kp.UserName())
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

func (me *KeyRing) fromPgpEntityList(lst ...*openpgp.Entity) *KeyRing {
	for _, entity := range lst {
		me.Add(keyFromEntity(entity))
	}
	return me
}

func (me *KeyRing) PubOnly() *KeyRing {
	ret := KeyRingCreate()
	for _, v := range me.kps {
		ret.Add(v.PubOnly())
	}
	return ret
}

func (me *KeyRing) ExportPubBinary(writer io.Writer) {
	for _, key := range me.kps {
		data := key.ExportPubBinary()
		writer.Write(data)
	}
}

func (me *KeyRing) ImportBinary(reader io.Reader) *KeyRing {
	lst, err := openpgp.ReadKeyRing(reader)
	util.Check(err)
	me.fromPgpEntityList(lst...)
	return me
}

func (me *KeyRing) ExportPubArmored(writer io.Writer) {
	enc := CreateEncrypter(writer).Armored(openpgp.PublicKeyType)
	defer enc.Close()
	me.ExportPubBinary(enc.Start())
}

func (me *KeyRing) ImportArmored(reader io.Reader) *KeyRing {
	lst, err := openpgp.ReadArmoredKeyRing(reader)
	util.Check(err)
	me.fromPgpEntityList(lst...)
	return me
}

func (me *KeyRing) ExportPrivBinary(writer io.Writer) {
	for _, key := range me.kps {
		data := key.ExportPrivBinary()
		writer.Write(data)
	}
}

func (me *KeyRing) ExportPrivArmored(writer io.Writer) {
	enc := CreateEncrypter(writer).Armored(openpgp.PublicKeyType)
	defer enc.Close()
	me.ExportPrivBinary(enc.Start())
}

type Abc struct {
	KeyRing
}

func (me *Abc) Ids() []uint64 {
	return []uint64{4}
}
