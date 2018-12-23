
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

func (vm *mockVM) canBranch(depth int) error {
	blockType := wasm.BlockTypeEmpty

	block := vm.getBlockFromDepth(depth)
	if block == nil {
		if depth == len(vm.blocks) {
			if len(vm.curFunc.ReturnTypes) != 0 {
				blockType = wasm.BlockType(vm.curFunc.ReturnTypes[0])
			}
		} else {
			return InvalidLabelError(uint32(depth))
		}
	} else if !block.loop {
		blockType = block.blockType
	}

	if blockType != wasm.BlockTypeEmpty {
		top, under := vm.topOperand()
		if under || top.Type != wasm.ValueType(blockType) {
			return InvalidTypeError{wasm.ValueType(blockType), top.Type}
		}
	}

	return nil
}


func (vm *mockVM) popBlock() *block {
	if len(vm.blocks) == 0 {
		return nil
	}

	stackTop := len(vm.blocks) - 1
	block := vm.blocks[stackTop]
	vm.blocks = append(vm.blocks[:stackTop], vm.blocks[stackTop+1:]...)

	return &block
}

func (vm *mockVM) topBlock() *block {
	if len(vm.blocks) == 0 {
		return nil
	}

	return &vm.blocks[len(vm.blocks)-1]
}

func (vm *mockVM) topOperand() (o operand, under bool) {
	stackTop := vm.stackTop - 1
	if stackTop == -1 {
		under = true
		return
	}
	o = vm.stack[stackTop]
	return
}

func (vm *mockVM) popOperand() (operand, bool) {
	var o operand
	stackTop := vm.stackTop - 1
	if stackTop == -1 {
		return o, true
	}
	o = vm.stack[stackTop]
	vm.stackTop--

	logger.Printf("Stack after pop is %v. Popped %v", vm.stack[:vm.stackTop], o)
	return o, false
}

func (vm *mockVM) pushOperand(t wasm.ValueType) {
	o := operand{t}
	logger.Printf("Stack top: %d, Len of stack :%d", vm.stackTop, len(vm.stack))
	if vm.stackTop == len(vm.stack) {
		vm.stack = append(vm.stack, o)
	} else {
		vm.stack[vm.stackTop] = o
	}
	vm.stackTop++

	logger.Printf("Stack after push is %v. Pushed %v", vm.stack[:vm.stackTop], o)
}

func (vm *mockVM) adjustStack(op ops.Op) error {
	for _, t := range op.Args {
		op, under := vm.popOperand()
		if !vm.isPolymorphic() && (under || op.Type != t) {
			return InvalidTypeError{t, op.Type}
		}
	}

	if op.Returns != wasm.ValueType(wasm.BlockTypeEmpty) {
		vm.pushOperand(op.Returns)
	}

	return nil
}

