package bloom

import (
	"math"

	"github.com/utkarsh-pro/use/pkg/structures/bitset"
)

// SBF represents a Standard Bloom filter.
type SBF struct {
	m uint64
	k uint64
	b *bitset.Bitset
}

// NewSBF returns a new pointer to a Standard Bloom filter.
func NewSBF(m, k uint64) *SBF {
	return &SBF{
		m: m,
		k: k,
		b: bitset.New(m),
	}
}

// Add adds an item to the Bloom filter.
func (f *SBF) Add(item []byte) {
	for _, i := range f.hash(item) {
		f.b.Set(i)
	}
}

// Contains returns true if the item is in the Bloom filter.
func (f *SBF) Contains(item []byte) bool {
	for _, i := range f.hash(item) {
		if !f.b.Get(i) {
			return false
		}
	}

	return true
}

// CurrentFalsePositiveRate returns the current false positive rate of the Bloom filter.
//
// Based on: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
func (f *SBF) CurrentFalsePositiveRate() float64 {
	// n is the number of items in the Bloom filter.
	n := f.ApproximateCount()

	// m is the total number of bits in the Bloom filter.
	m := f.b.Size()

	// k is the number of hash functions used.
	k := f.k

	return math.Pow((1 - math.Exp(float64((-1*int(k)*n)/int(m)))), float64(k))
}

// ApproximateCount returns the approximate number of items in the Bloom filter.
//
// Based on: https://en.wikipedia.org/wiki/Bloom_filter#Approximating_the_number_of_items_in_a_Bloom_filter
func (f *SBF) ApproximateCount() int {
	// X is the number of bits set to 1 in the Bloom filter.
	X := f.b.Count()

	// m is the total number of bits in the Bloom filter.
	m := f.b.Size()

	// k is the number of hash functions used.
	k := f.k

	return int((m / k) * uint64(math.Log(1-(float64(X)/float64(m)))))
}

// EstimateParameters estimates requirements for m and k.
//
// Based on: https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func EstimateParameters(n uint, e float64) (uint64, uint64) {
	m := uint64(math.Ceil(-1 * float64(n) * math.Log(e) / math.Pow(math.Log(2), 2)))
	k := uint64(math.Ceil(math.Log(2) * (float64(m) / float64(n))))
	return m, k
}

// NewSBFWithEstimates returns a pointer to a standard Bloom filter with the given number of items and
// false positive rate.
func NewSBFWithEstimates(n uint, e float64) *SBF {
	m, k := EstimateParameters(n, e)
	return NewSBF(m, k)
}

// Clear clears the Bloom filter.
func (f *SBF) Clear() {
	f.b.Clear()
}

// hash returns the hash values for the given item.
func (f *SBF) hash(item []byte) []uint64 {
	hashes := make([]uint64, f.k)

	for i := uint64(0); i < f.k; i++ {
		hashes[i] = f.hashFn(item, i)
	}

	return hashes
}

// hashFn returns the hash value for the given item.
func (f *SBF) hashFn(item []byte, i uint64) uint64 {
	var hash uint64

	for _, b := range item {
		// It's expected that this will overflow.
		hash = hash*31 + uint64(b)
	}

	return (hash + i) % f.m
}
