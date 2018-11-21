package vm

import (
	"XBlock/vm/errors"
	"XBlock/vm/types"
	"encoding/binary"
	"math/big"
	"reflect"
)

type BigIntSorter []big.Int

func (c BigIntSorter) Len() int {
	return len(c)
}
func (c BigIntSorter) Swap(i, j int) {
	if i >= 0 && i < len(c) && j >= 0 && j < len(c) { 		c[i], c[j] = c[j], c[i]
	}
}
