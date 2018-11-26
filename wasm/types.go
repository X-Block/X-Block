
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

func (t *ValueType) UnmarshalWASM(r io.Reader) error {
	v, err := leb128.ReadVarint32(r)
	if err != nil {
		return err
	}
	*t = ValueType(v)
	return nil
}

func (t ValueType) MarshalWASM(w io.Writer) error {
	_, err := leb128.WriteVarint64(w, int64(t))
	return err
}


type BlockType ValueType 
const BlockTypeEmpty BlockType = -0x40

func (b BlockType) String() string {
	if b == BlockTypeEmpty {
		return "<empty block>"
	}
	return ValueType(b).String()
}


type ElemType int 

const ElemTypeAnyFunc ElemType = -0x10

func (t *ElemType) UnmarshalWASM(r io.Reader) error {
	b, err := leb128.ReadVarint32(r)
	if err != nil {
		return err
	}
	*t = ElemType(b)
	return nil
}

func (t ElemType) MarshalWASM(w io.Writer) error {
	_, err := leb128.WriteVarint64(w, int64(t))
	return err
}

func (t ElemType) String() string {
	if t == ElemTypeAnyFunc {
		return "anyfunc"
	}

	return "<unknown elem_type>"
}


type FunctionSig struct {
	Form int8
	ParamTypes  []ValueType
	ReturnTypes []ValueType
}

