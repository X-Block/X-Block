
package exec

import (
	"math"
)


func (vm *VM) i32ReinterpretF32() {
	vm.pushUint32(math.Float32bits(vm.popFloat32()))
}

func (vm *VM) i64ReinterpretF64() {
	vm.pushUint64(math.Float64bits(vm.popFloat64()))
}

func (vm *VM) f32ReinterpretI32() {
	vm.pushFloat32(math.Float32frombits(vm.popUint32()))
}

func (vm *VM) f64ReinterpretI64() {
	vm.pushFloat64(math.Float64frombits(vm.popUint64()))
}
