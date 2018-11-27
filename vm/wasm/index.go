
package wasm

import (
	"fmt"
	"reflect"
)

type InvalidTableIndexError uint32

func (e InvalidTableIndexError) Error() string {
	return fmt.Sprintf("wasm: Invalid table to table index space: %d", uint32(e))
}

type InvalidValueTypeInitExprError struct {
	Wanted reflect.Kind
	Got    reflect.Kind
}

func (e InvalidValueTypeInitExprError) Error() string {
	return fmt.Sprintf("wasm: Wanted initializer expression to return %v value, got %v", e.Wanted, e.Got)
}

type InvalidLinearMemoryIndexError uint32

func (e InvalidLinearMemoryIndexError) Error() string {
	return fmt.Sprintf("wasm: Invalid linear memory index: %d", uint32(e))
}


func (m *Module) populateFunctions() error {
	if m.Types == nil || m.Function == nil {
		return nil
	}

	for codeIndex, typeIndex := range m.Function.Types {
		if int(typeIndex) >= len(m.Types.Entries) {
			return InvalidFunctionIndexError(typeIndex)
		}

		fn := Function{
			Sig:  &m.Types.Entries[typeIndex],
			Body: &m.Code.Bodies[codeIndex],
		}

		m.FunctionIndexSpace = append(m.FunctionIndexSpace, fn)
	}

	funcs := make([]uint32, 0, len(m.Function.Types)+len(m.imports.Funcs))
	funcs = append(funcs, m.imports.Funcs...)
	funcs = append(funcs, m.Function.Types...)
	m.Function.Types = funcs
	return nil
}

func (m *Module) GetFunction(i int) *Function {
	if i >= len(m.FunctionIndexSpace) || i < 0 {
		return nil
	}

	return &m.FunctionIndexSpace[i]
}

func (m *Module) populateGlobals() error {
	if m.Global == nil {
		return nil
	}

	m.GlobalIndexSpace = append(m.GlobalIndexSpace, m.Global.Globals...)
	logger.Printf("There are %d entries in the global index spaces.", len(m.GlobalIndexSpace))
	return nil
}

func (m *Module) GetGlobal(i int) *GlobalEntry {
	if i >= len(m.GlobalIndexSpace) || i < 0 {
		return nil
	}

	return &m.GlobalIndexSpace[i]
}

func (m *Module) populateTables() error {
	if m.Table == nil || len(m.Table.Entries) == 0 || m.Elements == nil || len(m.Elements.Entries) == 0 {
		return nil
	}

	for _, elem := range m.Elements.Entries {

		if int(elem.Index) >= len(m.TableIndexSpace) {
			return InvalidTableIndexError(elem.Index)
		}

		val, err := m.ExecInitExpr(elem.Offset)
		if err != nil {
			return err
		}
		offset, ok := val.(int32)
		if !ok {
			return InvalidValueTypeInitExprError{reflect.Int32, reflect.TypeOf(val).Kind()}
		}

		table := m.TableIndexSpace[int(elem.Index)]
		if int(offset)+len(elem.Elems) > len(table) {
			data := make([]uint32, int(offset)+len(elem.Elems))
			copy(data[offset:], elem.Elems)
			copy(data, table)
			m.TableIndexSpace[int(elem.Index)] = data
		} else {
			copy(table[int(offset):], elem.Elems)
			m.TableIndexSpace[int(elem.Index)] = table
		}
	}

	logger.Printf("There are %d entries in the table index space.", len(m.TableIndexSpace))
	return nil
}

func (m *Module) GetTableElement(index int) (uint32, error) {
	if index >= len(m.TableIndexSpace[0]) {
		return 0, InvalidTableIndexError(index)
	}

	return m.TableIndexSpace[0][index], nil
}

func (m *Module) populateLinearMemory() error {
	if m.Data == nil || len(m.Data.Entries) == 0 {
		return nil
	}
	
	for _, entry := range m.Data.Entries {
		if entry.Index != 0 {
			return InvalidLinearMemoryIndexError(entry.Index)
		}

		val, err := m.ExecInitExpr(entry.Offset)
		if err != nil {
			return err
		}
		offset, ok := val.(int32)
		if !ok {
			return InvalidValueTypeInitExprError{reflect.Int32, reflect.TypeOf(val).Kind()}
		}

		memory := m.LinearMemoryIndexSpace[int(entry.Index)]
		if int(offset)+len(entry.Data) > len(memory) {
			data := make([]byte, int(offset)+len(entry.Data))
			copy(data[offset:], entry.Data)
			copy(data, memory)
			m.LinearMemoryIndexSpace[int(entry.Index)] = data
		} else {
			copy(memory[int(offset):], entry.Data)
			m.LinearMemoryIndexSpace[int(entry.Index)] = memory
		}
	}

	return nil
}

func (m *Module) GetLinearMemoryData(index int) (byte, error) {
	if index >= len(m.LinearMemoryIndexSpace[0]) {
		return 0, InvalidLinearMemoryIndexError(uint32(index))

	}

	return m.LinearMemoryIndexSpace[0][index], nil
}
