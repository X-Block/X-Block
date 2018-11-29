
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


func (m *Module) Custom(name string) *SectionCustom {
	for _, s := range m.Customs {
		if s.Name == name {
			return s
		}
	}
	return nil
}


func NewModule() *Module {
	return &Module{
		Types:    &SectionTypes{},
		Import:   &SectionImports{},
		Table:    &SectionTables{},
		Memory:   &SectionMemories{},
		Global:   &SectionGlobals{},
		Export:   &SectionExports{},
		Start:    &SectionStartFunction{},
		Elements: &SectionElements{},
		Data:     &SectionData{},
	}
}

type ResolveFunc func(name string) (*Module, error)

func DecodeModule(r io.Reader) (*Module, error) {
	reader := &readpos.ReadPos{
		R:      r,
		CurPos: 0,
	}
	m := &Module{}
	magic, err := readU32(reader)
	if err != nil {
		return nil, err
	}
	if magic != Magic {
		return nil, ErrInvalidMagic
	}
	if m.Version, err = readU32(reader); err != nil {
		return nil, err
	}

	for {
		done, err := m.readSection(reader)
		if err != nil {
			return nil, err
		} else if done {
			return m, nil
		}
	}
}

