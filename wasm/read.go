
package wasm

import (
	"encoding/binary"
	"io"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

func readBytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := io.ReadFull(r, bytes)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func readBytesUint(r io.Reader) ([]byte, error) {
	n, err := leb128.ReadVarUint32(r)
	if err != nil {
		return nil, err
	}
	return readBytes(r, int(n))
}

func readString(r io.Reader, n int) (string, error) {
	bytes, err := readBytes(r, n)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

