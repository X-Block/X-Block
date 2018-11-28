
package wasm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/X-Block/X-Block/tree/master/vm/wasm/leb128"
)

const (
	i32Const  byte = 0x41
	i64Const  byte = 0x42
	f32Const  byte = 0x43
	f64Const  byte = 0x44
	getGlobal byte = 0x23
	end       byte = 0x0b
)

var ErrEmptyInitExpr = errors.New("wasm: Initializer expression produces no value")

type InvalidInitExprOpError byte

func (e InvalidInitExprOpError) Error() string {
	return fmt.Sprintf("wasm: Invalid opcode in initializer expression: %#x", byte(e))
}

type InvalidGlobalIndexError uint32

func (e InvalidGlobalIndexError) Error() string {
	return fmt.Sprintf("wasm: Invalid index to global index space: %#x", uint32(e))
}

func readInitExpr(r io.Reader) ([]byte, error) {
	b := make([]byte, 1)
	buf := new(bytes.Buffer)
	r = io.TeeReader(r, buf)

outer:
	for {
		_, err := io.ReadFull(r, b)
		if err != nil {
			return nil, err
		}
		switch b[0] {
		case i32Const:
			_, err := leb128.ReadVarint32(r)
			if err != nil {
				return nil, err
			}
		case i64Const:
			_, err := leb128.ReadVarint64(r)
			if err != nil {
				return nil, err
			}
		case f32Const:
			if _, err := readU32(r); err != nil {
				return nil, err
			}
		case f64Const:
			if _, err := readU64(r); err != nil {
				return nil, err
			}
		case getGlobal:
			_, err := leb128.ReadVarUint32(r)
			if err != nil {
				return nil, err
			}
		case end:
			break outer
		default:
			return nil, InvalidInitExprOpError(b[0])
		}
	}

	if buf.Len() == 0 {
		return nil, ErrEmptyInitExpr
	}

	return buf.Bytes(), nil
}
