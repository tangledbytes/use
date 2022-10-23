package tlvrw

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/utkarsh-pro/use/pkg/utils"
)

var (
	MaxLen uint32 = math.MaxUint32 // 4GB max at once

	ErrInvalidLen    = fmt.Errorf("length cannot be greater than %d", MaxLen)
	ErrInvalidWhence = fmt.Errorf("invalid whence")
)

type Reader struct {
	// r is the reader which will be used to read from the
	// underlying resource.
	r io.ReaderAt

	// readerPos is the current position of the reader.
	readerPos int64
}

type Writer struct {
	// w is the writer which will be used to write to the
	// underlying resource.
	w io.Writer
}

func NewReader(r io.ReaderAt) *Reader {
	return &Reader{
		r:         r,
		readerPos: 0,
	}
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

// Read reads a TLV from the reader.
func (r *Reader) Read(tlv *TLV) error {
	if err := r.ReadLazy(tlv); err != nil {
		return err
	}

	return r.Fill(tlv)
}

// ReadLazy reads a TLV from the reader, but does not read the value.
// This is useful for when you want to read the type and length, but
// not the value.
//
// The value will be read when the Fill(tlv) method is called.
func (r *Reader) ReadLazy(tlv *TLV) error {
	// Read 1 byte for the type.
	typBytes := make([]byte, 1)
	if n, err := r.r.ReadAt(typBytes, r.readerPos); err != nil {
		if err == io.EOF {
			// If we get an EOF, we need to check if we read anything.
			if n == 0 {
				return io.EOF
			}

			// If we read something, we need to return an error.
			return io.ErrUnexpectedEOF
		}

		return err
	} else {
		r.readerPos += int64(n)
		tlv.Typ = typBytes[0]
	}

	// Read 4 bytes for the length.
	lenBytes := make([]byte, 4)
	if n, err := r.r.ReadAt(lenBytes, r.readerPos); err != nil {
		if err == io.EOF {
			// regardless of what we read, we need to return an error
			// as we are in middle of reading a TLV
			return io.ErrUnexpectedEOF
		}

		return err
	} else {
		r.readerPos += int64(n)
		tlv.Len = binary.LittleEndian.Uint32(lenBytes)
	}

	// Skip reading value but store the position for future reads
	tlv.Val = nil
	tlv.valuepos = utils.ToPointer(r.readerPos)

	// Move the reader position to the end of the value
	r.readerPos += int64(tlv.Len)

	return nil
}

// Fill takes a partially filled TLV and fills in the value. This function
// can be called as many times as needed (in case of failure to read earlier).
func (r *Reader) Fill(tlv *TLV) error {
	if tlv.Done() {
		return nil
	}

	if tlv.Len > MaxLen {
		return ErrInvalidLen
	}

	if tlv.Val == nil {
		tlv.Val = make([]byte, tlv.Len)
	}

	if n, err := r.r.ReadAt(tlv.Val, *tlv.valuepos); err != nil {
		if err == io.EOF {
			if n < int(tlv.Len) {
				return io.ErrUnexpectedEOF
			}

			return io.EOF
		}

		return err
	}

	// Clear the valuepos only if we are able to read the value
	// successfully. This is to ensure that we can read the value
	// again if it fails once.
	tlv.valuepos = nil

	return nil
}

// Seek seeks to the given offset, supported whence are io.SeekStart and
// io.SeekCurrent. io.SeekEnd is not supported.
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.readerPos = offset
	case io.SeekCurrent:
		r.readerPos += offset
	default:
		return 0, ErrInvalidWhence
	}

	return r.readerPos, nil
}

// Write writes a TLV to the writer.
func (w *Writer) Write(tlv *TLV) error {
	if tlv.Len > MaxLen {
		return ErrInvalidLen
	}

	if err := binary.Write(w.w, binary.LittleEndian, tlv.Typ); err != nil {
		return err
	}

	if err := binary.Write(w.w, binary.LittleEndian, tlv.Len); err != nil {
		return err
	}

	for {
		n, err := w.w.Write(tlv.Val)
		if err != nil {
			return err
		}

		if n == len(tlv.Val) {
			break
		}
	}

	return nil
}
