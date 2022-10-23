package stupid

import (
	"io"

	"github.com/utkarsh-pro/use/pkg/storage/tlvrw"
)

const (
	// OpTypeTLV is the type of the operation TLV.
	OpTypeTLV = byte(1)

	// KeyTypeTLV is the type of the key TLV.
	KeyTypeTLV = byte(2)

	// ValTypeTLV is the type of the value TLV.
	ValTypeTLV = byte(3)
)

type Packet struct {
	Op  byte
	Key []byte
	Val []byte

	vtlv *tlvrw.TLV
}

type PacketLite struct {
	Op     byte
	Key    []byte
	ValPos int64
}

type reader struct {
	// r is the underlying TLV reader
	r *tlvrw.Reader
}

type writer struct {
	// w is the underlying TLV writer
	w *tlvrw.Writer
}

func newreader(r io.ReaderAt) *reader {
	return &reader{
		r: tlvrw.NewReader(r),
	}
}

func newwriter(w io.Writer) *writer {
	return &writer{
		w: tlvrw.NewWriter(w),
	}
}

// lread is a lazy reader which reads the packet
func (r *reader) lread(p *Packet) error {
	// read the operation type TLV
	optlv := tlvrw.NewTLV(OpTypeTLV, nil)
	if err := r.r.Read(optlv); err != nil {
		// EOF indicates that there are no more packets to read
		if err == io.EOF {
			return io.EOF
		}

		return err
	}
	p.Op = optlv.Val[0]

	// read the key type TLV
	keytlv := tlvrw.NewTLV(KeyTypeTLV, nil)
	if err := r.r.Read(keytlv); err != nil {
		if err == io.EOF {
			// packets are set of 3 TLVs and we don't expect
			// EOF on the second TLV read
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

func (w *writer) write(p *Packet) error {
	// Write an operation type TLV
	if err := w.w.Write(tlvrw.NewTLV(OpTypeTLV, []byte{p.Op})); err != nil {
		return err
	}

	// Write a key type TLV
	if err := w.w.Write(tlvrw.NewTLV(KeyTypeTLV, p.Key)); err != nil {
		return err
	}

	// Write a value type TLV
	if err := w.w.Write(tlvrw.NewTLV(ValTypeTLV, p.Val)); err != nil {
		return err
	}

	return nil
}
