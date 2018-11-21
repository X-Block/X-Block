package types

import (
	"math/big"
	"XBlock/vm/interfaces"
)

type InteropInterface struct {
	_object interfaces.IInteropInterface
}

func NewInteropInterface(value interfaces.IInteropInterface) *InteropInterface {
	var ii InteropInterface
	ii._object = value
	return &ii
}}
