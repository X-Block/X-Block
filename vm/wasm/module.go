
package wasm

import (
	"errors"
	"io"
	"reflect"

	"github.com/X-Block/X-Block/tree/master/vm/wasm/internal/readpos"
)

var ErrInvalidMagic = errors.New("wasm: Invalid magic number")

type Function struct {
	Sig  *FunctionSig
	Body *FunctionBody
	Host reflect.Value
}

func (fct *Function) IsHost() bool {
	return fct.Host != reflect.Value{}
}

type Module struct {
	Version  uint32
	Sections []Section

	Types    *SectionTypes
	Import   *SectionImports
	Function *SectionFunctions
	Table    *SectionTables
	Memory   *SectionMemories
	Global   *SectionGlobals
	Export   *SectionExports
	Start    *SectionStartFunction
	Elements *SectionElements
	Code     *SectionCode
	Data     *SectionData
	Customs  []*SectionCustom


	FunctionIndexSpace []Function
	GlobalIndexSpace   []GlobalEntry

	TableIndexSpace        [][]uint32
	LinearMemoryIndexSpace [][]byte

	imports struct {
		Funcs    []uint32
		Globals  int
		Tables   int
		Memories int
	}
}


