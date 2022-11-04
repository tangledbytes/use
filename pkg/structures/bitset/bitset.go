package bitset

// BitsetUnitSize is the size of a bitset unit.
const BitsetUnitSize = 64

type Bitset struct {
	// bits is the underlying bitset.
	bits []uint64
	// size is the size of the bitset.
	size uint64
}

// New returns a new Bitset instance.
func New(size uint64) *Bitset {
	return &Bitset{
		bits: make([]uint64, size/BitsetUnitSize+1),
		size: size,
	}
}

// Set sets the bit at the given index.
func (b *Bitset) Set(index uint64) {
	if index >= b.size {
		return
	}

	// A | B
	b.bits[index/BitsetUnitSize] |= 1 << (index % BitsetUnitSize)
}

// Unset unsets the bit at the given index.
func (b *Bitset) Unset(index uint64) {
	if index >= b.size {
		return
	}

	// A & ~B
	b.bits[index/BitsetUnitSize] &= ^(1 << (index % BitsetUnitSize)) // &^ === AND NOT
}

// Get returns the bit at the given index.
func (b *Bitset) Get(index uint64) bool {
	if index >= b.size {
		return false
	}

	// A & B
	return b.bits[index/BitsetUnitSize]&(1<<(index%BitsetUnitSize)) != 0
}

// Size returns the size of the bitset.
func (b *Bitset) Size() uint64 {
	return b.size
}

// Bits returns the underlying bitset.
func (b *Bitset) Bits() []uint64 {
	return b.bits
}

// Clear clears the bitset.
func (b *Bitset) Clear() {
	for i := range b.bits {
		b.bits[i] = 0
	}
}

// Copy returns a copy of the bitset.
func (b *Bitset) Copy() *Bitset {
	return &Bitset{
		bits: append([]uint64{}, b.bits...),
		size: b.size,
	}
}

// Union returns the union of the bitset with the given bitset.
func (b *Bitset) Union(other *Bitset) *Bitset {
	if b.size != other.size {
		return nil
	}

	result := New(b.size)
	for i := range b.bits {
		result.bits[i] = b.bits[i] | other.bits[i]
	}

	return result
}

// Intersect returns the intersection of the bitset with the given bitset.
func (b *Bitset) Intersect(other *Bitset) *Bitset {
	if b.size != other.size {
		return nil
	}

	result := New(b.size)
	for i := range b.bits {
		result.bits[i] = b.bits[i] & other.bits[i]
	}

	return result
}

// Difference returns the difference of the bitset with the given bitset.
func (b *Bitset) Difference(other *Bitset) *Bitset {
	if b.size != other.size {
		return nil
	}

	result := New(b.size)
	for i := range b.bits {
		result.bits[i] = b.bits[i] &^ other.bits[i]
	}

	return result
}

// Count returns the number of 1s in the bitset.
func (b *Bitset) Count() int {
	var count int
	for _, bits := range b.bits {
		for bits != 0 {
			bits &= bits - 1
			count++
		}
	}

	return count
}
