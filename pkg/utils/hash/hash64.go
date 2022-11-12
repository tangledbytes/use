package hash

import (
	gohash "hash"
	"io"
)

// DoubleHash64 satisfies the hash.Hash64 interface.
//
// It is a double hash function that uses two hash functions to generate
// a hash value.
//
// Note:
//   - `Sum` is a noop.
//   - `Size` is a noop
//   - `BlockSize` is a noop.
type DoubleHash64 struct {
	io.Writer

	fn1 gohash.Hash64
	fn2 gohash.Hash64
	i   uint64
}

func NewDoubleHash64(fn1, fn2 gohash.Hash64, i uint64) gohash.Hash64 {
	return &DoubleHash64{
		fn1: fn1,
		fn2: fn2,
		i:   i,
	}
}

func DoubleHash64KGenerator(fn1, fn2 gohash.Hash64, k uint64) []gohash.Hash64 {
	hashes := make([]gohash.Hash64, k)

	for i := uint64(0); i < k; i++ {
		hashes[i] = NewDoubleHash64(fn1, fn2, i)
	}

	return hashes
}

func (h *DoubleHash64) Write(p []byte) (n int, err error) {
	h.fn1.Write(p)
	h.fn2.Write(p)

	return len(p), nil
}

func (h *DoubleHash64) Reset() {
	h.fn1.Reset()
	h.fn2.Reset()
}

func (h *DoubleHash64) Sum(b []byte) []byte {
	// noop
	return b
}

func (h *DoubleHash64) Size() int {
	// noop
	return 0
}

func (h *DoubleHash64) BlockSize() int {
	// noop
	return 0
}

func (h *DoubleHash64) Sum64() uint64 {
	return h.fn1.Sum64() + h.i*h.fn2.Sum64() + (h.i * h.i)
}
