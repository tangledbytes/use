package hash

import (
	gohash "hash"
	"io"
)

// DoubleHash32 satisfies the hash.Hash32 interface.
//
// It is a double hash function that uses two hash functions to generate
// a hash value.
//
// Note:
//   - `Sum` is a noop.
//   - `Size` is a noop
//   - `BlockSize` is a noop.
type DoubleHash32 struct {
	io.Writer

	fn1 gohash.Hash32
	fn2 gohash.Hash32
	i   uint32
}

func NewDoubleHash32(fn1, fn2 gohash.Hash32, i uint32) gohash.Hash32 {
	return &DoubleHash32{
		fn1: fn1,
		fn2: fn2,
		i:   i,
	}
}

func DoubleHash32KGenerator(fn1, fn2 gohash.Hash32, k uint32) []gohash.Hash32 {
	hashes := make([]gohash.Hash32, k)

	for i := uint32(0); i < k; i++ {
		hashes[i] = NewDoubleHash32(fn1, fn2, i)
	}

	return hashes
}

func (h *DoubleHash32) Write(p []byte) (n int, err error) {
	h.fn1.Write(p)
	h.fn2.Write(p)

	return len(p), nil
}

func (h *DoubleHash32) Reset() {
	h.fn1.Reset()
	h.fn2.Reset()
}

func (h *DoubleHash32) Sum(b []byte) []byte {
	// noop
	return b
}

func (h *DoubleHash32) Size() int {
	// noop
	return 0
}

func (h *DoubleHash32) BlockSize() int {
	// noop
	return 0
}

func (h *DoubleHash32) Sum32() uint32 {
	return h.fn1.Sum32() + h.i*h.fn2.Sum32() + (h.i * h.i)
}
