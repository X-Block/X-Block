
import (
	"bytes"
	"io"

)


func verifyBody(fn *wasm.FunctionSig, body *wasm.FunctionBody, module *wasm.Module) (*mockVM, error) {
	vm := &mockVM{
		stack:    []operand{},
		stackTop: 0,

		code:       bytes.NewReader(body.Code),
		origLength: len(body.Code),

		polymorphic: false,
		blocks:      []block{},
		curFunc:     fn,
	}

	localVariables := []operand{}


	for _, entry := range fn.ParamTypes {
		localVariables = append(localVariables, operand{entry})
	}

	for _, entry := range body.Locals {
		vars := make([]operand, entry.Count)
		for i := uint32(0); i < entry.Count; i++ {
			vars[i].Type = entry.Type
			logger.Printf("Var %v", entry.Type)
		}
		localVariables = append(localVariables, vars...)
	}

	for {
		op, err := vm.code.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return vm, err
		}

		opStruct, err := ops.New(op)
		if err != nil {
			return vm, err
		}

		logger.Printf("PC: %d OP: %s polymorphic: %v", vm.pc(), opStruct.Name, vm.isPolymorphic())

		if !opStruct.Polymorphic {
			if err := vm.adjustStack(opStruct); err != nil {
				return vm, err
			}
		}

		switch op {
		case ops.If, ops.Block, ops.Loop:
			sig, err := vm.fetchVarInt()
			if err != nil {
				return vm, err
			}

			switch wasm.ValueType(sig) {
			case wasm.ValueTypeI32, wasm.ValueTypeI64, wasm.ValueTypeF32, wasm.ValueTypeF64, wasm.ValueType(wasm.BlockTypeEmpty):
				vm.pushBlock(op, wasm.BlockType(sig))
			default:
				if !vm.isPolymorphic() {
					return vm, InvalidImmediateError{"block_type", opStruct.Name}
				}
			}

		case ops.Else:
			block := vm.topBlock()
			if block == nil || block.op != ops.If {
				return vm, UnmatchedOpError(op)
			}

			if block.blockType != wasm.BlockTypeEmpty {
				top, under := vm.topOperand()
				if !vm.isPolymorphic() && (under || top.Type != wasm.ValueType(block.blockType)) {
					return vm, InvalidTypeError{wasm.ValueType(block.blockType), top.Type}
				}
				vm.pushOperand(wasm.ValueType(block.blockType))
			}
			vm.stackTop = block.stackTop
		case ops.End:
			isPolymorphic := vm.isPolymorphic()

			block := vm.popBlock()
			if block == nil {
				return vm, UnmatchedOpError(op)
			}

			if block.blockType != wasm.BlockTypeEmpty {
				top, under := vm.topOperand()
				if !isPolymorphic && (under || top.Type != wasm.ValueType(block.blockType)) {
					return vm, InvalidTypeError{wasm.ValueType(block.blockType), top.Type}
				}
				vm.stackTop = block.stackTop
				vm.pushOperand(wasm.ValueType(block.blockType))
				vm.stackTop = block.stackTop + 1 
			} else {
				vm.stackTop = block.stackTop
			}

		case ops.BrIf, ops.Br:
			depth, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			if err = vm.canBranch(int(depth)); !vm.isPolymorphic() && err != nil {
				return vm, err
			}
			if op == ops.Br {
				vm.setPolymorphic()
			}
		case ops.BrTable:
			operand, under := vm.popOperand()
			if !vm.isPolymorphic() && (under || operand.Type != wasm.ValueTypeI32) {
				return vm, InvalidTypeError{wasm.ValueTypeI32, operand.Type}
			}
		
			targetCount, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			var targetTable []uint32
			for i := uint32(0); i < targetCount; i++ {
				entry, err := vm.fetchVarUint()
				if err != nil {
					return vm, err
				}
				if err = vm.canBranch(int(entry)); !vm.isPolymorphic() && err != nil {
					return vm, err
				}
				targetTable = append(targetTable, entry)
			}

			defaultTarget, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			if err = vm.canBranch(int(defaultTarget)); !vm.isPolymorphic() && err != nil {
				return vm, err
			}
			vm.setPolymorphic()

		case ops.Return:
			if len(fn.ReturnTypes) > 1 {
				panic("not implemented")
			}
			if len(fn.ReturnTypes) != 0 {
				if !vm.isPolymorphic() && (under || top.Type != fn.ReturnTypes[0]) {
					return vm, InvalidTypeError{fn.ReturnTypes[0], top.Type}
				}
			}
			vm.setPolymorphic()

		case ops.Unreachable:
			vm.setPolymorphic()

		case ops.I32Const:
			_, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
		case ops.I64Const:
			_, err := vm.fetchVarInt64()
			if err != nil {
				return vm, err
			}
		case ops.F32Const:
			_, err := vm.fetchUint32()
			if err != nil {
				return vm, err
			}
		case ops.F64Const:
			_, err := vm.fetchUint64()
			if err != nil {
				return vm, err
			}
		case ops.GetLocal, ops.SetLocal, ops.TeeLocal:
			i, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			if int(i) >= len(localVariables) {
				return vm, InvalidLocalIndexError(i)
			}

			v := localVariables[i]

			if op == ops.GetLocal {
				vm.pushOperand(v.Type)
			} else { 
                   if !vm.isPolymorphic() && (under || top.Type != v.Type) {
					return vm, InvalidTypeError{v.Type, top.Type}
				}
				if op == ops.TeeLocal {
					vm.pushOperand(v.Type)
				}
			}

		case ops.GetGlobal, ops.SetGlobal:
			index, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			gv := module.GetGlobal(int(index))
			if gv == nil {
				return vm, wasm.InvalidGlobalIndexError(index)
			}
			if op == ops.GetGlobal {
				vm.pushOperand(gv.Type.Type)
			} else {
				val, under := vm.popOperand()
				if !vm.isPolymorphic() && (under || val.Type != gv.Type.Type) {
					return vm, InvalidTypeError{gv.Type.Type, val.Type}
				}
			}

		case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			
			_, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			
			_, err = vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
		case ops.CurrentMemory, ops.GrowMemory:
			_, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

		case ops.Call:
			index, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			fn := module.GetFunction(int(index))
			if fn == nil {
				return vm, wasm.InvalidFunctionIndexError(index)
			}

			logger.Printf("Function being called: %v", fn)
			for index := range fn.Sig.ParamTypes {
				argType := fn.Sig.ParamTypes[len(fn.Sig.ParamTypes)-index-1]
				operand, under := vm.popOperand()
				if !vm.isPolymorphic() && (under || operand.Type != argType) {
					return vm, InvalidTypeError{argType, operand.Type}
				}
			}

			if len(fn.Sig.ReturnTypes) > 0 {
				vm.pushOperand(fn.Sig.ReturnTypes[0])
			}

		case ops.CallIndirect:
			if module.Table == nil || len(module.Table.Entries) == 0 {
				return vm, NoSectionError(wasm.SectionIDTable)
			}
					index, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			fnExpectSig := module.Types.Entries[index]

			if operand, under := vm.popOperand(); !vm.isPolymorphic() && (under || operand.Type != wasm.ValueTypeI32) {
				return vm, InvalidTypeError{wasm.ValueTypeI32, operand.Type}
			}

			for index := range fnExpectSig.ParamTypes {
				argType := fnExpectSig.ParamTypes[len(fnExpectSig.ParamTypes)-index-1]
				operand, under := vm.popOperand()
				if !vm.isPolymorphic() && (under || (operand.Type != argType)) {
					return vm, InvalidTypeError{argType, operand.Type}
				}
			}

			if len(fnExpectSig.ReturnTypes) > 0 {
				vm.pushOperand(fnExpectSig.ReturnTypes[0])
			}

		case ops.Drop:
			if _, under := vm.popOperand(); !vm.isPolymorphic() && under {
				return vm, ErrStackUnderflow
			}

		case ops.Select:
			if vm.isPolymorphic() {
				continue
			}
			operands := make([]operand, 2)
			c, under := vm.popOperand()
			if under || c.Type != wasm.ValueTypeI32 {
				return vm, InvalidTypeError{wasm.ValueTypeI32, c.Type}
			}

			for i := 0; i < 2; i++ {
				operand, under := vm.popOperand()
				if !vm.isPolymorphic() && under {
					return vm, ErrStackUnderflow
				}
				operands[i] = operand
			}

			if operands[0].Type != operands[1].Type {
				return vm, InvalidTypeError{operands[1].Type, operands[2].Type}
			}

			vm.pushOperand(operands[1].Type)
		}
	}

	return vm, nil
}

func VerifyModule(module *wasm.Module) error {
	if module.Function == nil || module.Types == nil || len(module.Types.Entries) == 0 {
		return nil
	}
	if module.Code == nil {
		return NoSectionError(wasm.SectionIDCode)
	}

	logger.Printf("There are %d functions", len(module.Function.Types))
	for i, fn := range module.FunctionIndexSpace {
		if vm, err := verifyBody(fn.Sig, fn.Body, module); err != nil {
			return Error{vm.pc(), i, err}
		}
		logger.Printf("No errors in function %d", i)
	}

	return nil
}
