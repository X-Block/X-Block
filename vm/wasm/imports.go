
package wasm

import (
	"errors"
	"fmt"
	"io"

	"github.com/X-Block/X-Block/tree/master/vm/wasm/leb128"
)

type Import interface {
	Kind() External
	Marshaler
	isImport()
}


type ImportEntry struct {
	ModuleName string 
	FieldName  string 

	Type Import
}

type FuncImport struct {
	Type uint32
}

func (FuncImport) isImport() {}
func (FuncImport) Kind() External {
	return ExternalFunction
}
func (f FuncImport) MarshalWASM(w io.Writer) error {
	_, err := leb128.WriteVarUint32(w, uint32(f.Type))
	return err
}

type TableImport struct {
	Type Table
}

func (TableImport) isImport() {}
func (TableImport) Kind() External {
	return ExternalTable
}
func (t TableImport) MarshalWASM(w io.Writer) error {
	return t.Type.MarshalWASM(w)
}

type MemoryImport struct {
	Type Memory
}

func (MemoryImport) isImport() {}
func (MemoryImport) Kind() External {
	return ExternalMemory
}
func (t MemoryImport) MarshalWASM(w io.Writer) error {
	return t.Type.MarshalWASM(w)
}

type GlobalVarImport struct {
	Type GlobalVar
}

