
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

func (GlobalVarImport) isImport() {}
func (GlobalVarImport) Kind() External {
	return ExternalGlobal
}
func (t GlobalVarImport) MarshalWASM(w io.Writer) error {
	return t.Type.MarshalWASM(w)
}

var (
	ErrImportMutGlobal           = errors.New("wasm: cannot import global mutable variable")
	ErrNoExportsInImportedModule = errors.New("wasm: imported module has no exports")
)

type InvalidExternalError uint8

func (e InvalidExternalError) Error() string {
	return fmt.Sprintf("wasm: invalid external_kind value %d", uint8(e))
}

type ExportNotFoundError struct {
	ModuleName string
	FieldName  string
}

type KindMismatchError struct {
	ModuleName string
	FieldName  string
	Import     External
	Export     External
}

func (e KindMismatchError) Error() string {
	return fmt.Sprintf("wasm: Mismatching import and export external kind values for %s.%s (%v, %v)", e.FieldName, e.ModuleName, e.Import, e.Export)
}

func (e ExportNotFoundError) Error() string {
	return fmt.Sprintf("wasm: couldn't find export with name %s in module %s", e.FieldName, e.ModuleName)
}

type InvalidFunctionIndexError uint32

func (e InvalidFunctionIndexError) Error() string {
	return fmt.Sprintf("wasm: Invalid index to function index space: %#x", uint32(e))
}

