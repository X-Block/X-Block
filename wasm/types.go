
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

func (f FunctionSig) String() string {
	return fmt.Sprintf("<func %v -> %v>", f.ParamTypes, f.ReturnTypes)
}

type InvalidTypeConstructorError struct {
	Wanted int
	Got    int
}

func (e InvalidTypeConstructorError) Error() string {
	return fmt.Sprintf("wasm: invalid type constructor: wanted %d, got %d", e.Wanted, e.Got)
}

func (f *FunctionSig) UnmarshalWASM(r io.Reader) error {
	form, err := leb128.ReadVarint32(r)
	if err != nil {
		return err
	}
	f.Form = int8(form)

	paramCount, err := leb128.ReadVarUint32(r)
	if err != nil {
		return err
	}
	f.ParamTypes = make([]ValueType, paramCount)

	for i := range f.ParamTypes {
		err = f.ParamTypes[i].UnmarshalWASM(r)
		if err != nil {
			return err
		}
	}

	returnCount, err := leb128.ReadVarUint32(r)
	if err != nil {
		return err
	}

	f.ReturnTypes = make([]ValueType, returnCount)
	for i := range f.ReturnTypes {
		err = f.ReturnTypes[i].UnmarshalWASM(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FunctionSig) MarshalWASM(w io.Writer) error {
	_, err := leb128.WriteVarint64(w, int64(f.Form))
	if err != nil {
		return err
	}

	_, err = leb128.WriteVarUint32(w, uint32(len(f.ParamTypes)))
	if err != nil {
		return err
	}
	for _, p := range f.ParamTypes {
		err = p.MarshalWASM(w)
		if err != nil {
			return err
		}
	}

	_, err = leb128.WriteVarUint32(w, uint32(len(f.ReturnTypes)))
	if err != nil {
		return err
	}
	for _, p := range f.ReturnTypes {
		err = p.MarshalWASM(w)
		if err != nil {
			return err
		}
	}
	return nil
}


type GlobalVar struct {
	Type    ValueType 
	Mutable bool      
}

func (g *GlobalVar) UnmarshalWASM(r io.Reader) error {
	*g = GlobalVar{}

	err := g.Type.UnmarshalWASM(r)
	if err != nil {
		return err
	}

	m, err := leb128.ReadVarUint32(r)
	if err != nil {
		return err
	}

	g.Mutable = m == 1

	return nil
}

func (g *GlobalVar) MarshalWASM(w io.Writer) error {
	if err := g.Type.MarshalWASM(w); err != nil {
		return err
	}
	var m uint32
	if g.Mutable {
		m = 1
	}
	if _, err := leb128.WriteVarUint32(w, m); err != nil {
		return err
	}
	return nil
}


type Table struct {
	ElementType ElemType
	Limits      ResizableLimits
}

func (t *Table) UnmarshalWASM(r io.Reader) error {
	err := t.ElementType.UnmarshalWASM(r)
	if err != nil {
		return err
	}

	err = t.Limits.UnmarshalWASM(r)
	if err != nil {
		return err
	}
	return err
}

func (t *Table) MarshalWASM(w io.Writer) error {
	if err := t.ElementType.MarshalWASM(w); err != nil {
		return err
	}
	if err := t.Limits.MarshalWASM(w); err != nil {
		return err
	}
	return nil
}

type Memory struct {
	Limits ResizableLimits
}

func (m *Memory) UnmarshalWASM(r io.Reader) error {
	return m.Limits.UnmarshalWASM(r)
}

func (m *Memory) MarshalWASM(w io.Writer) error {
	return m.Limits.MarshalWASM(w)
}

type External uint8

const (
	ExternalFunction External = 0
	ExternalTable    External = 1
	ExternalMemory   External = 2
	ExternalGlobal   External = 3
)

func (e External) String() string {
	switch e {
	case ExternalFunction:
		return "function"
	case ExternalTable:
		return "table"
	case ExternalMemory:
		return "memory"
	case ExternalGlobal:
		return "global"
	default:
		return "<unknown external_kind>"
	}
}
func (e *External) UnmarshalWASM(r io.Reader) error {
	bytes, err := readBytes(r, 1)
	if err != nil {
		return err
	}
	*e = External(bytes[0])
	return nil
}
func (e External) MarshalWASM(w io.Writer) error {
	_, err := w.Write([]byte{byte(e)})
	return err
}

type ResizableLimits struct {
	Flags   uint32 
	Initial uint32 
	Maximum uint32 
}

func (lim *ResizableLimits) UnmarshalWASM(r io.Reader) error {
	*lim = ResizableLimits{}
	f, err := leb128.ReadVarUint32(r)
	if err != nil {
		return err
	}
	lim.Flags = f

	lim.Initial, err = leb128.ReadVarUint32(r)
	if err != nil {
		return err
	}

	if lim.Flags&0x1 != 0 {
		m, err := leb128.ReadVarUint32(r)
		if err != nil {
			return err
		}
		lim.Maximum = m
	}
	return nil
}

func (lim *ResizableLimits) MarshalWASM(w io.Writer) error {
	if _, err := leb128.WriteVarUint32(w, uint32(lim.Flags)); err != nil {
		return err
	}
	if _, err := leb128.WriteVarUint32(w, uint32(lim.Initial)); err != nil {
		return err
	}
	if lim.Flags&0x1 != 0 {
		if _, err := leb128.WriteVarUint32(w, uint32(lim.Maximum)); err != nil {
			return err
		}
	}
	return nil
}
