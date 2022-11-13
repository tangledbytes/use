package types

import (
	"hash/maphash"

	"github.com/utkarsh-pro/use/pkg/hash"
)

// Hash type is a hash function.
//
// Hash functions should consume a byte slice and return a uint64.
type Hash func([]byte) uint64

// DefaultHash returns a slice of k default hash functions.
func DefaultHash(k uint64) []Hash {
	hashFns := make([]Hash, k)

	for i := uint64(0); i < k; i++ {
		var h1 maphash.Hash
		h1.SetSeed(maphash.MakeSeed())

		var h2 maphash.Hash
		h2.SetSeed(maphash.MakeSeed())

		h := hash.NewDoubleHash64(&h1, &h2, i)

		hashFns[i] = func(b []byte) uint64 {
			h.Write(b)
			hashed := h.Sum64()

			h.Reset()
			return hashed
		}
	}

	return hashFns
}
