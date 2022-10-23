package stupid

import "encoding/binary"

func encodeid(id uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, id)
	return b
}

func decodeid(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}
