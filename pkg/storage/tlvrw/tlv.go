package tlvrw

import "fmt"

// TLV represents a type,length,value structure.
type TLV struct {
	Typ byte
	Len uint32
	Val []byte

	valuepos *int64
}

// NewTLV returns a pointer to a new TLV.
func NewTLV(typ byte, val []byte) *TLV {
	return &TLV{
		Typ: typ,
		Len: uint32(len(val)),
		Val: val,
	}
}

// Done returns true if the TLV has been read completely.
//
// This is useful for lazy reads where TLV just has the type and
// length, but not the value.
func (tlv *TLV) Done() bool {
	return tlv.valuepos == nil && tlv.Len == uint32(len(tlv.Val))
}

// ValuePos returns the position of the value bytes that is saved
// when the TLV is read lazily.
func (tlv *TLV) ValuePos() int64 {
	return *tlv.valuepos
}

// String returns a string representation of the TLV.
func (tlv *TLV) String() string {
	return fmt.Sprintf("TLV{Typ: %d, Len: %d, Val: %v}", tlv.Typ, tlv.Len, tlv.Val)
}
