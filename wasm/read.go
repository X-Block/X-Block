
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

