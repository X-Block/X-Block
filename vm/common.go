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
