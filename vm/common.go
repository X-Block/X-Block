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

func MaxBigInt(ints []big.Int) big.Int{
	max := ints[0]

	for _, d := range ints {
		if d.Cmp(&max) > 0 {
			max = d
		}
	}

	return max
}

func MinInt64(datas []int64) int64 {

	var minimum int64
	for i, d := range datas { 
		if i == 0 {
			minimum = d
		}
		if d < minimum {
			minimum = d
		}
	}

	return minimum
}

func MaxInt64(datas []int64) int64 {

	var maximum int64
	
	for i, d := range datas { 
		if i == 0 {
			maximum = d
			
		}
		if d > maximum {
			maximum = d
		}
	}

	return maximum
}

func Concat(array1 []byte, array2 []byte) []byte {
	len := len(array2)
	for i := 0; i < len; i++ {
		array1 = append(array1, array2[i]) 
	}

	return array1
}

func BigIntOp(bi *big.Int, op OpCode) *big.Int {
	var nb *big.Int
	switch op {
	case INC:
		nb = bi.Add(bi, big.NewInt(int64(1)))
	case DEC:
		nb = bi.Sub(bi, big.NewInt(int64(1)))
	case SAL:
		nb = bi.Lsh(bi, 1)
	case SAR:
		nb = bi.Rsh(bi, 1)
	case NEGATE:
		nb = bi.Neg(bi)
	case ABS:
		nb = bi.Abs(bi)
	default:
		nb = bi
	}
	return nb
}

func AsBool(e interface{}) bool {
	if v, ok := e.([]byte); ok {
		for _, b := range v {
			if b != 0 {
				return true
			}
		}
	}
	return false
}

