
package wasm

import (
	"fmt"
	"io"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

type Marshaler interface {

	MarshalWASM(w io.Writer) error
}

type Unmarshaler interface {

	UnmarshalWASM(r io.Reader) error
}


type ValueType int8

const (
	ValueTypeI32 ValueType = -0x01
	ValueTypeI64 ValueType = -0x02
	ValueTypeF32 ValueType = -0x03
	ValueTypeF64 ValueType = -0x04
)

var valueTypeStrMap = map[ValueType]string{
	ValueTypeI32: "i32",
	ValueTypeI64: "i64",
	ValueTypeF32: "f32",
	ValueTypeF64: "f64",
}

func (t ValueType) String() string {
	str, ok := valueTypeStrMap[t]
	if !ok {
		str = fmt.Sprintf("<unknown value_type %d>", int8(t))
	}
	return str
}


const TypeFunc int = -0x20

