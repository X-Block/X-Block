
package validate

import (
	"bytes"
	"encoding/binary"
	"io"
)

type mockVM struct {
	stack      []operand
	stackTop   int 
	origLength int 

	code *bytes.Reader

	polymorphic bool   
	blocks      []block 

	curFunc *wasm.FunctionSig
}

type block struct {
	pc          int            
	stackTop    int            
	blockType   wasm.BlockType 
	op          byte           
	polymorphic bool           
	loop        bool           
}

func (vm *mockVM) fetchVarUint() (uint32, error) {
	return leb128.ReadVarUint32(vm.code)
}

func (vm *mockVM) fetchVarInt() (int32, error) {
	return leb128.ReadVarint32(vm.code)
}

func (vm *mockVM) fetchVarInt64() (int64, error) {
	return leb128.ReadVarint64(vm.code)
}

func (vm *mockVM) fetchUint32() (uint32, error) {
	var buf [4]byte
	_, err := io.ReadFull(vm.code, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func (vm *mockVM) fetchUint64() (uint64, error) {
	var buf [8]byte
	_, err := io.ReadFull(vm.code, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

func (vm *mockVM) pushBlock(op byte, blockType wasm.BlockType) {
	logger.Printf("Pushing block %v", blockType)
	vm.blocks = append(vm.blocks, block{
		pc:          vm.pc(),
		stackTop:    vm.stackTop,
		blockType:   blockType,
		polymorphic: vm.isPolymorphic(),
		op:          op,
		loop:        op == ops.Loop,
	})
}


func (vm *mockVM) getBlockFromDepth(depth int) *block {
	if depth >= len(vm.blocks) {
		return nil
	}

	return &vm.blocks[len(vm.blocks)-1-depth]
}

