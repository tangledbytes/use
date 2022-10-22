package tlvrw

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

var (
	MaxLen uint32 = math.MaxUint32 // 4GB max at once

	ErrInvalidLen    = fmt.Errorf("length cannot be greater than %d", MaxLen)
	EOF              = io.EOF
	ErrUnexpectedEOF = fmt.Errorf("unexpected EOF")
)

// TLV represents a type,length,value structure.
type TLV struct {
	Typ byte
	Len uint32
	Val []byte
}

// NewTLV returns a pointer to a new TLV.
func NewTLV(typ byte, val []byte) *TLV {
	return &TLV{
		Typ: typ,
		Len: uint32(len(val)),
		Val: val,
	}
}

// ResetReader resets the reader to the beginning of the file.
func ResetReader(r io.ReadSeeker) error {
	_, err := r.Seek(0, io.SeekStart)
	return err
}

// Read reads a TLV from the reader.
func Read(r io.Reader) (*TLV, error) {
	tlv := NewTLV(0, nil)

	if err := binary.Read(r, binary.LittleEndian, &tlv.Typ); err != nil {
		if err == io.EOF {
			// return nil as nothing was really read
			return nil, EOF
		}

		return tlv, err
	}

	if err := binary.Read(r, binary.LittleEndian, &tlv.Len); err != nil {
		if err == io.EOF {
			// It is not expected that the data will end here.
			return tlv, ErrUnexpectedEOF
		}

		return tlv, err
	}

	tlv.Val = make([]byte, tlv.Len)
	if n, err := r.Read(tlv.Val); err != nil {
		if err == io.EOF {
			if n == 0 {
				// It is not expected that the data will end here.
				return tlv, ErrUnexpectedEOF
			}

			return tlv, EOF
		}

		return tlv, err
	}

	return tlv, nil
}

// Skip skips the current TLV.
func Skip(r io.ReadSeeker) error {
	tlv := NewTLV(0, nil)

	if err := binary.Read(r, binary.LittleEndian, &tlv.Typ); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &tlv.Len); err != nil {
		return err
	}

	if _, err := r.Seek(int64(tlv.Len), io.SeekCurrent); err != nil {
		return err
	}

	return nil
}

// Write writes a TLV to the io.Writer.
func Write(w io.Writer, tlv *TLV) error {
	if tlv.Len > MaxLen {
		return ErrInvalidLen
	}

	if err := binary.Write(w, binary.LittleEndian, tlv.Typ); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, tlv.Len); err != nil {
		return err
	}

	for {
		n, err := w.Write(tlv.Val)
		if err != nil {
			return err
		}

		if n == len(tlv.Val) {
			break
		}
	}

	return nil
}
