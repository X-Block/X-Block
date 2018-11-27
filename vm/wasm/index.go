
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

