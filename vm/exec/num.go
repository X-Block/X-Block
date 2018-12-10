
package exec

import (
	"math"
	"math/bits"
)


func (vm *VM) i32Clz() {
	vm.pushUint64(uint64(bits.LeadingZeros32(vm.popUint32())))
}

func (vm *VM) i32Ctz() {
	vm.pushUint64(uint64(bits.TrailingZeros32(vm.popUint32())))
}

func (vm *VM) i32Popcnt() {
	vm.pushUint64(uint64(bits.OnesCount32(vm.popUint32())))
}

func (vm *VM) i32Add() {
	vm.pushUint32(vm.popUint32() + vm.popUint32())
}

func (vm *VM) i32Mul() {
	vm.pushUint32(vm.popUint32() * vm.popUint32())
}

func (vm *VM) i32DivS() {
	v2 := vm.popInt32()
	v1 := vm.popInt32()
	vm.pushInt32(v1 / v2)
}

func (vm *VM) i32DivU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(v1 / v2)
}

func (vm *VM) i32RemS() {
	v2 := vm.popInt32()
	v1 := vm.popInt32()
	vm.pushInt32(v1 % v2)
}

func (vm *VM) i32RemU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(v1 % v2)
}

func (vm *VM) i32Sub() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(v1 - v2)
}

func (vm *VM) i32And() {
	vm.pushUint32(vm.popUint32() & vm.popUint32())
}

func (vm *VM) i32Or() {
	vm.pushUint32(vm.popUint32() | vm.popUint32())
}

func (vm *VM) i32Xor() {
	vm.pushUint32(vm.popUint32() ^ vm.popUint32())
}

func (vm *VM) i32Shl() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(v1 << v2)
}

func (vm *VM) i32ShrU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(v1 >> v2)
}

func (vm *VM) i32ShrS() {
	v2 := vm.popUint32()
	v1 := vm.popInt32()
	vm.pushInt32(v1 >> v2)
}

func (vm *VM) i32Rotl() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(bits.RotateLeft32(v1, int(v2)))
}

func (vm *VM) i32Rotr() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushUint32(bits.RotateLeft32(v1, -int(v2)))
}

func (vm *VM) i32LeS() {
	v2 := vm.popInt32()
	v1 := vm.popInt32()
	vm.pushBool(v1 <= v2)
}

func (vm *VM) i32LeU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushBool(v1 <= v2)
}

func (vm *VM) i32LtS() {
	v2 := vm.popInt32()
	v1 := vm.popInt32()
	vm.pushBool(v1 < v2)
}

func (vm *VM) i32LtU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushBool(v1 < v2)
}

func (vm *VM) i32GtS() {
	v2 := vm.popInt32()
	v1 := vm.popInt32()
	vm.pushBool(v1 > v2)
}

func (vm *VM) i32GeS() {
	v2 := vm.popInt32()
	v1 := vm.popInt32()
	vm.pushBool(v1 >= v2)
}

func (vm *VM) i32GtU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushBool(v1 > v2)
}

func (vm *VM) i32GeU() {
	v2 := vm.popUint32()
	v1 := vm.popUint32()
	vm.pushBool(v1 >= v2)
}

func (vm *VM) i32Eqz() {
	vm.pushBool(vm.popUint32() == 0)
}

func (vm *VM) i32Eq() {
	vm.pushBool(vm.popUint32() == vm.popUint32())
}

func (vm *VM) i32Ne() {
	vm.pushBool(vm.popUint32() != vm.popUint32())
}


func (vm *VM) i64Clz() {
	vm.pushUint64(uint64(bits.LeadingZeros64(vm.popUint64())))
}

func (vm *VM) i64Ctz() {
	vm.pushUint64(uint64(bits.TrailingZeros64(vm.popUint64())))
}

func (vm *VM) i64Popcnt() {
	vm.pushUint64(uint64(bits.OnesCount64(vm.popUint64())))
}

func (vm *VM) i64Add() {
	vm.pushUint64(vm.popUint64() + vm.popUint64())
}

func (vm *VM) i64Sub() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushUint64(v1 - v2)
}

func (vm *VM) i64Mul() {
	vm.pushUint64(vm.popUint64() * vm.popUint64())
}

func (vm *VM) i64DivS() {
	v2 := vm.popInt64()
	v1 := vm.popInt64()
	vm.pushInt64(v1 / v2)
}

func (vm *VM) i64DivU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushUint64(v1 / v2)
}

func (vm *VM) i64RemS() {
	v2 := vm.popInt64()
	v1 := vm.popInt64()
	vm.pushInt64(v1 % v2)
}

func (vm *VM) i64RemU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushUint64(v1 % v2)
}

func (vm *VM) i64And() {
	vm.pushUint64(vm.popUint64() & vm.popUint64())
}

func (vm *VM) i64Or() {
	vm.pushUint64(vm.popUint64() | vm.popUint64())
}

func (vm *VM) i64Xor() {
	vm.pushUint64(vm.popUint64() ^ vm.popUint64())
}

func (vm *VM) i64Shl() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushUint64(v1 << v2)
}

func (vm *VM) i64ShrS() {
	v2 := vm.popUint64()
	v1 := vm.popInt64()
	vm.pushInt64(v1 >> v2)
}

func (vm *VM) i64ShrU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushUint64(v1 >> v2)
}

func (vm *VM) i64Rotl() {
	v2 := vm.popInt64()
	v1 := vm.popUint64()
	vm.pushUint64(bits.RotateLeft64(v1, int(v2)))
}

func (vm *VM) i64Rotr() {
	v2 := vm.popInt64()
	v1 := vm.popUint64()
	vm.pushUint64(bits.RotateLeft64(v1, -int(v2)))
}

func (vm *VM) i64Eq() {
	vm.pushBool(vm.popUint64() == vm.popUint64())
}

func (vm *VM) i64Eqz() {
	vm.pushBool(vm.popUint64() == 0)
}

func (vm *VM) i64Ne() {
	vm.pushBool(vm.popUint64() != vm.popUint64())
}

func (vm *VM) i64LtS() {
	v2 := vm.popInt64()
	v1 := vm.popInt64()
	vm.pushBool(v1 < v2)
}

func (vm *VM) i64LtU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushBool(v1 < v2)
}

func (vm *VM) i64GtS() {
	v2 := vm.popInt64()
	v1 := vm.popInt64()
	vm.pushBool(v1 > v2)
}

func (vm *VM) i64GtU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushBool(v1 > v2)
}

func (vm *VM) i64LeU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushBool(v1 <= v2)
}

func (vm *VM) i64LeS() {
	v2 := vm.popInt64()
	v1 := vm.popInt64()
	vm.pushBool(v1 <= v2)
}

func (vm *VM) i64GeS() {
	v2 := vm.popInt64()
	v1 := vm.popInt64()
	vm.pushBool(v1 >= v2)
}

func (vm *VM) i64GeU() {
	v2 := vm.popUint64()
	v1 := vm.popUint64()
	vm.pushBool(v1 >= v2)
}


