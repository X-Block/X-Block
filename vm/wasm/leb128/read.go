package leb128

import (
	"io"
)

func ReadVarUint32Size(r io.Reader) (res uint32, size uint, err error) {
	b := make([]byte, 1)
	var shift uint
	for {
		if _, err = io.ReadFull(r, b); err != nil {
			return
		}

		size++

		cur := uint32(b[0])
		res |= (cur & 0x7f) << (shift)
		if cur&0x80 == 0 {
			return res, size, nil
		}
		shift += 7
	}
}

func ReadVarUint32(r io.Reader) (uint32, error) {
	n, _, err := ReadVarUint32Size(r)
	return n, err
}

