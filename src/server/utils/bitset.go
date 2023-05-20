package utils

import "github.com/bits-and-blooms/bitset"

type AuthBitSet struct {
	*bitset.BitSet
}

func NewBitSet(length uint, auths ...uint) AuthBitSet {
	var auth AuthBitSet
	auth.BitSet = bitset.New(length)
	for _, i := range auths {
		auth.Set(i)
	}
	return auth
}
