package exec

func (vm *VM) i32Const() {
	vm.pushUint32(vm.fetchUint32())
}

func (vm *VM) i64Const() {
	vm.pushUint64(vm.fetchUint64())
}

func (vm *VM) f32Const() {
	vm.pushFloat32(vm.fetchFloat32())
}

func (vm *VM) f64Const() {
	vm.pushFloat64(vm.fetchFloat64())
}
