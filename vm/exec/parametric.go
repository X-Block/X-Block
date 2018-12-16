
package exec

func (vm *VM) drop() {
	vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-1]
}

func (vm *VM) selectOp() {
	c := vm.popUint32()
	val2 := vm.popUint64()
	val1 := vm.popUint64()

	if c != 0 {
		vm.pushUint64(val1)
	} else {
		vm.pushUint64(val2)
	}
}
