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

func AsInt64(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if len(b) > 8 {
		return 0, errors.ErrBadValue
	}

	var bs [8]byte
	copy(bs[:], b)

	res := binary.LittleEndian.Uint64(bs[:])

	return int64(res), nil
}

func ByteArrZip(s1 []byte, s2 []byte, op OpCode) []byte{
	var ns []byte
	switch op {
	case CAT:
		ns = append(s1, s2...)
	}
	return ns
}

func BigIntZip(ints1 *big.Int, ints2 *big.Int, op OpCode) *big.Int {
	var nb *big.Int
	switch op {
	case AND:
		nb = ints1.And(ints1, ints2)
	case OR:
		nb = ints1.Or(ints1, ints2)
	case XOR:
		nb = ints1.Xor(ints1, ints2)
	case ADD:
		nb = ints1.Add(ints1, ints2)
	case SUB:
		nb = ints1.Sub(ints1, ints2)
	case MUL:
		nb = ints1.Mul(ints1, ints2)
	case DIV:
		nb = ints1.Div(ints1, ints2)
	case MOD:
		nb = ints1.Mod(ints1, ints2)
	case SHL:
		nb = ints1.Lsh(ints1, uint(ints2.Int64()))
	case SHR:
		nb = ints1.Rsh(ints1, uint(ints2.Int64()))
	case MIN:
		c := ints1.Cmp(ints2)
		if c <= 0 {
			nb = ints1
		} else {
			nb = ints2
		}
	case MAX:
		c := ints1.Cmp(ints2)
		if c <= 0 {
			nb = ints2
		} else {
			nb = ints1
		}
	}
	return nb
}

func BigIntComp(bigint *big.Int, op OpCode) bool {
	var nb bool
	switch op {
	case NZ:
		nb = bigint.Cmp(big.NewInt(int64(0))) != 0
	}
	return nb
}

