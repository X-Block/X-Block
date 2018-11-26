
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

