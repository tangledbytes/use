package stupid

import (
	"bytes"
	"io"

	"github.com/utkarsh-pro/use/pkg/tlvrw"
)

const (
	// IDTypeTLV is the type of the ID TLV
	IDTypeTLV = byte(1)

	// OpTypeTLV is the type of the operation TLV.
	OpTypeTLV = byte(2)

	// KeyTypeTLV is the type of the key TLV.
	KeyTypeTLV = byte(3)

	// ValTypeTLV is the type of the value TLV.
	ValTypeTLV = byte(4)
)

// Packet is the basic unit of data transfer between the engine
// and the file system.
type Packet struct {
	ID  uint64
	Op  byte
	Key []byte
	Val []byte

	vtlv *tlvrw.TLV
}

// reader is a packet reader
type reader struct {
	// r is the underlying TLV reader
	r *tlvrw.Reader
}

// writer is a packet writer
type writer struct {
	// w is the underlying TLV writer
	w io.Writer
}

// newreader returns a new packet reader
func newreader(r io.ReaderAt) *reader {
	return &reader{
		r: tlvrw.NewReader(r),
	}
}

// newwriter returns a new packet writer
func newwriter(w io.Writer) *writer {
	return &writer{
		w: w,
	}
}

// lread is a lazy reader which reads the packet
func (r *reader) lread(p *Packet) error {
	// read the ID type TLV
	idtlv := tlvrw.NewTLV(IDTypeTLV, nil)
	if err := r.r.Read(idtlv); err != nil {
		// EOF indicates that there are no more packets to read
		if err == io.EOF {
			return io.EOF
		}

		return err
	}
	p.ID = decodeid(idtlv.Val)

	// read the operation type TLV
	optlv := tlvrw.NewTLV(OpTypeTLV, nil)
	if err := r.r.Read(optlv); err != nil {
		if err == io.EOF {
			// packets are set of 4 TLVs and we don't expect
			// EOF on the second TLV read
			return io.ErrUnexpectedEOF
		}

		return err
	}
	p.Op = optlv.Val[0]

	// read the key type TLV
	keytlv := tlvrw.NewTLV(KeyTypeTLV, nil)
	if err := r.r.Read(keytlv); err != nil {
		if err == io.EOF {
			// packets are set of 4 TLVs and we don't expect
			// EOF on the third TLV read
			return io.ErrUnexpectedEOF
		}

		return err
	}
	p.Key = keytlv.Val

	// read the value type TLV lazily
	valtlv := tlvrw.NewTLV(ValTypeTLV, nil)
	if err := r.r.ReadLazy(valtlv); err != nil {
		return err
	}
	p.vtlv = valtlv

	return nil
}

// fill fills the value type TLV with the value
func (r *reader) fill(p *Packet) error {
	if p.vtlv == nil {
		return nil
	}

	if err := r.r.Fill(p.vtlv); err != nil {
		return err
	}

	p.Val = p.vtlv.Val
	p.vtlv = nil
	return nil
}

// pos returns the current position of the reader
func (r *reader) pos() int64 {
	pos, _ := r.r.Seek(0, io.SeekCurrent)
	return pos
}

// write writes the packet to the underlying writer
//
// write will copy the packet to a buffer and write the buffer
// to the underlying writer in one go. This is to ensure that the
// packet is written in one go and is not broken in between.
func (w *writer) write(p *Packet) error {
	// get estimate size of the packet
	sizeOfIDInBytes := uint32(8) // 8 bytes for uint64
	sizeOfOpInBytes := uint32(1) // 1 byte for byte
	sizeOfKeyInBytes := uint32(len(p.Key))
	sizeOfValInBytes := uint32(len(p.Val))

	size := sizeOfIDInBytes + sizeOfOpInBytes + sizeOfKeyInBytes + sizeOfValInBytes

	// buffer data and write it in one go
	buf := bytes.NewBuffer(make([]byte, 0, size))

	// create a TLV writer
	tw := tlvrw.NewWriter(buf)

	// write the ID type TLV
	if err := tw.Write(tlvrw.NewTLV(IDTypeTLV, encodeid(p.ID))); err != nil {
		return err
	}

	// Write an operation type TLV
	if err := tw.Write(tlvrw.NewTLV(OpTypeTLV, []byte{p.Op})); err != nil {
		return err
	}

	// Write a key type TLV
	if err := tw.Write(tlvrw.NewTLV(KeyTypeTLV, p.Key)); err != nil {
		return err
	}

	// Write a value type TLV
	if err := tw.Write(tlvrw.NewTLV(ValTypeTLV, p.Val)); err != nil {
		return err
	}

	// write the buffer to the underlying writer
	if _, err := w.w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
