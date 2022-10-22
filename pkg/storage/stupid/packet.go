package stupid

import (
	"fmt"
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

var (
	ErrUnexpectedEOF = fmt.Errorf("unexpected EOF")
)

type Packet struct {
	Op  byte
	Key []byte
	Val []byte
}

type PacketLite struct {
	Op     byte
	Key    []byte
	ValPos int64
}

// ReadPacket reads a packet from the reader.
func ReadPacket(r io.ReadSeeker) (*Packet, error) {
	// Read an operation type TLV
	op, err := readOpType(r)
	if err != nil {
		if err == tlvrw.EOF {
			if op == 0 {
				return nil, err
			}

			return nil, ErrUnexpectedEOF
		}

		return nil, err
	}

	// Read a key type TLV
	key, err := readKey(r)
	if err != nil {
		if err == tlvrw.EOF {
			return nil, ErrUnexpectedEOF
		}

		return nil, err
	}

	// Read a value type TLV
	val, err := readValue(r)
	if err != nil {
		return nil, err
	}

	return &Packet{
		Op:  op,
		Key: key,
		Val: val,
	}, nil
}

// ReadPacketLite reads the packet lite from the reader.
func ReadPacketLite(r io.ReadSeeker) (*PacketLite, error) {
	// Read an operation type TLV
	op, err := readOpType(r)
	if err != nil {
		if err == tlvrw.EOF {
			if op == 0 {
				return nil, err
			}

			return nil, ErrUnexpectedEOF
		}

		return nil, err
	}

	// Read a key type TLV
	key, err := readKey(r)
	if err != nil {
		if err == tlvrw.EOF {
			return nil, ErrUnexpectedEOF
		}

		return nil, err
	}

	// Save the position of the value type TLV
	valSeekPos, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	if err := tlvrw.Skip(r); err != nil {
		return nil, err
	}

	return &PacketLite{
		Op:     op,
		Key:    key,
		ValPos: valSeekPos,
	}, nil
}

func ResetReader(r io.ReadSeeker) error {
	return tlvrw.ResetReader(r)
}

// WritePacket writes a packet to the writer.
func WritePacket(w io.Writer, p *Packet) error {
	// Write an operation type TLV
	if err := tlvrw.Write(w, tlvrw.NewTLV(OpTypeTLV, []byte{p.Op})); err != nil {
		return err
	}

	// Write a key type TLV
	if err := tlvrw.Write(w, tlvrw.NewTLV(KeyTypeTLV, p.Key)); err != nil {
		return err
	}

	// Write a value type TLV
	if err := tlvrw.Write(w, tlvrw.NewTLV(ValTypeTLV, p.Val)); err != nil {
		return err
	}

	return nil
}

func readOpType(r io.ReadSeeker) (byte, error) {
	optlv, err := tlvrw.Read(r)
	if err != nil && err != tlvrw.EOF {
		return 0, err
	}

	if optlv == nil && err == tlvrw.EOF {
		return 0, err
	}

	if optlv.Typ != OpTypeTLV {
		return 0, fmt.Errorf("expected operation type TLV, got %d", optlv.Typ)
	}
	return optlv.Val[0], err
}

func readKey(r io.ReadSeeker) ([]byte, error) {
	keytlv, err := tlvrw.Read(r)
	if err != nil && err != tlvrw.EOF {
		return nil, err
	}
	if keytlv.Typ != KeyTypeTLV {
		return nil, fmt.Errorf("expected key type TLV, got %d", keytlv.Typ)
	}
	return keytlv.Val, nil
}

func readValue(r io.ReadSeeker) ([]byte, error) {
	valtlv, err := tlvrw.Read(r)
	if err != nil && err != tlvrw.EOF {
		return nil, err
	}
	if valtlv.Typ != ValTypeTLV {
		return nil, fmt.Errorf("expected value type TLV, got %d", valtlv.Typ)
	}
	return valtlv.Val, nil
}

func readValueAt(r io.ReadSeeker, pos int64) ([]byte, error) {
	if _, err := r.Seek(pos, io.SeekStart); err != nil {
		return nil, err
	}

	return readValue(r)
}
