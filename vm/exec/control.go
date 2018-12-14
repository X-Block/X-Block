
package exec

import "errors"

var ErrUnreachable = errors.New("exec: reached unreachable")

func (vm *VM) unreachable() {
	panic(ErrUnreachable)
}

func (vm *VM) nop() {}
