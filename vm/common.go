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
func (c BigIntSorter) Less(i, j int) bool {
	if i >= 0 && i < len(c) && j >= 0 && j < len(c) { 		return c[i].Cmp(&c[j]) < 0
	}

	return false
}

func ToBigInt(data interface{}) *big.Int {
	var bi big.Int
	switch t := data.(type) {
	case int64:
		bi.SetInt64(int64(t))
	case int32:
		bi.SetInt64(int64(t))
	case int16:
		bi.SetInt64(int64(t))
	case int8:
		bi.SetInt64(int64(t))
	case int:
		bi.SetInt64(int64(t))
	case uint64:
		bi.SetUint64(uint64(t))
	case uint32:
		bi.SetUint64(uint64(t))
	case uint16:
		bi.SetUint64(uint64(t))
	case uint8:
		bi.SetUint64(uint64(t))
	case uint:
		bi.SetUint64(uint64(t))
	case big.Int:
		bi = t
	case *big.Int:
		bi = *t
	}
	return &bi
}

//common func
func SumBigInt(ints []big.Int) big.Int {
	sum := big.NewInt(0)
	for _, v := range ints {
		sum = sum.Add(sum, &v)
	}
	return *sum
}

func MinBigInt(ints []big.Int) big.Int{
	minimum := ints[0]

	for _, d := range ints {
		if d.Cmp(&minimum) < 0 {
			minimum = d
		}
	}

	return minimum
}

