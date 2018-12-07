
package exec

import "errors"

var (
	ErrSignatureMismatch = errors.New("exec: signature mismatch in call_indirect")
	ErrUndefinedElementIndex = errors.New("exec: undefined element index")
)

func (vm *VM) call() {
	index := vm.fetchUint32()

	vm.funcs[index].call(vm, int64(index))
}

func (vm *VM) callIndirect() {
	index := vm.fetchUint32()
	fnExpect := vm.module.Types.Entries[index]
	_ = vm.fetchUint32() 
	tableIndex := vm.popUint32()
	if int(tableIndex) >= len(vm.module.TableIndexSpace[0]) {
		panic(ErrUndefinedElementIndex)
	}
	elemIndex := vm.module.TableIndexSpace[0][tableIndex]
	fnActual := vm.module.FunctionIndexSpace[elemIndex]

	if len(fnExpect.ParamTypes) != len(fnActual.Sig.ParamTypes) {
		panic(ErrSignatureMismatch)
	}
	if len(fnExpect.ReturnTypes) != len(fnActual.Sig.ReturnTypes) {
		panic(ErrSignatureMismatch)
	}

	for i := range fnExpect.ParamTypes {
		if fnExpect.ParamTypes[i] != fnActual.Sig.ParamTypes[i] {
			panic(ErrSignatureMismatch)
		}
	}

	for i := range fnExpect.ReturnTypes {
		if fnExpect.ReturnTypes[i] != fnActual.Sig.ReturnTypes[i] {
			panic(ErrSignatureMismatch)
		}
	}

	vm.funcs[elemIndex].call(vm, int64(elemIndex))
}
