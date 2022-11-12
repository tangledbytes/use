package dibf

import (
	"math"

	"github.com/utkarsh-pro/use/pkg/structures/bitset"
	"github.com/utkarsh-pro/use/pkg/structures/bloom/types"
	"github.com/utkarsh-pro/use/pkg/utils"
)

// DIBF represents a Deletable Bloom filter.
//
// Based on: https://arxiv.org/pdf/1005.0352.pdf
type DIBF struct {
	m uint64
	k uint64
	r *bitset.Bitset
	b *bitset.Bitset

	hashFns []types.Hash
}

// New returns a new pointer to a Deletable Bloom filter.
//
//	m - the total number of bits in the Bloom filter.
//	k - the number of hash functions used.
//	r - the size of the bitset used to store the collsions (r <= m).
func New(m, k, r uint64, hashFns []types.Hash) *DIBF {
	if r == 0 || r > m {
		r = m
	}

	return &DIBF{
		m:       utils.Max(m, 1),
		k:       utils.Max(k, 1),
		r:       bitset.New(r),
		b:       bitset.New(utils.Max(m, 1)),
		hashFns: hashFns,
	}
}

// EstimateParameters estimates requirements for m and k.
//
// Based on: https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func EstimateParameters(n uint, e float64) (uint64, uint64) {
	m := uint64(math.Ceil(-1 * (float64(n) * math.Log(e)) / math.Pow(math.Log(2), 2)))
	k := uint64(math.Ceil(-1 * math.Log(e) / math.Log(2)))
	return m, k
}

// NewWithEstimates returns a pointer to a standard Bloom filter with the given number of items and
// false positive rate.
func NewWithEstimates(n uint, e float64, f float64, hashFns []types.Hash) *DIBF {
	m, k := EstimateParameters(n, e)
	return New(m, k, uint64(float64(m)/f), hashFns)
}

// Add adds an item to the Bloom filter.
func (f *DIBF) Add(item []byte) {
	for _, i := range f.hash(item) {
		// if the bit is already set, mark the collision bit.
		if f.b.Get(i) {
			f.r.Set(f.getRegion(i))
		}

		f.b.Set(i)
	}
}

// Contains returns true if the item is in the Bloom filter.
func (f *DIBF) Contains(item []byte) bool {
	for _, i := range f.hash(item) {
		if !f.b.Get(i) {
			return false
		}
	}

	return true
}

// Delete deletes an item from the Bloom filter.
func (f *DIBF) Delete(item []byte) {
	hashes := f.hash(item)

	// if the item is not in the Bloom filter then return.
	if !f.Contains(item) {
		return
	}

	// if the item is in the Bloom filter then try to delete it.
	for _, i := range hashes {
		// if the collision bit is set then we cannot delete the item.
		if f.r.Get(f.getRegion(i)) {
			return
		}

		f.b.Unset(i)
	}
}

// CurrentFalsePositiveRate returns the current false positive rate of the Bloom filter.
//
// Based on: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
func (f *DIBF) CurrentFalsePositiveRate() float64 {
	// n is the number of items in the Bloom filter.
	n := f.ApproximateCount()

	// m is the total number of bits in the Bloom filter.
	m := f.b.Size()

	// k is the number of hash functions used.
	k := f.k

	return math.Pow((1 - math.Exp((float64(-1*int(k)*n) / float64(m)))), float64(k))
}

// ApproximateCount returns the approximate number of items in the Bloom filter.
//
// Based on: https://en.wikipedia.org/wiki/Bloom_filter#Approximating_the_number_of_items_in_a_Bloom_filter
func (f *DIBF) ApproximateCount() int {
	// X is the number of bits set to 1 in the Bloom filter.
	X := float64(f.b.Count())

	// m is the total number of bits in the Bloom filter.
	m := float64(f.b.Size())

	// k is the number of hash functions used.
	k := float64(f.k)

	count := (m / k) * math.Log(1/(1-X/m))
	return int(math.Floor(count + 0.5))
}

// K returns the number of hash functions used.
func (f *DIBF) K() uint64 {
	return f.k
}

// Cap returns the capacity of the Bloom filter.
func (f *DIBF) Cap() uint64 {
	return f.m
}

// Clear clears the Bloom filter.
func (f *DIBF) Clear() {
	f.b.Clear()
}

func (f *DIBF) getRegion(i uint64) uint64 {
	regionSize := f.m / f.r.Size()
	return i / regionSize
}

// hash returns the hash values for the given item.
func (f *DIBF) hash(item []byte) []uint64 {
	hashes := make([]uint64, f.k)

	// if no hash functions were provided then use the default hash functions.
	if f.hashFns == nil {
		f.hashFns = types.DefaultHash(f.k)
	}

	for i := uint64(0); i < f.k; i++ {
		hashes[i] = f.hashFns[i](item) % f.m
	}

	return hashes
}
