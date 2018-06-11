/*
 * Copyright (C) 2018 The X-Block Authors
 * This file is part of The X-Block library.
 *
 * The X-Block is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The X-Block is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The X-Block.  If not, see <http://www.gnu.org/licenses/>.
 */

package serialize

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var (
	// ErrWrongUintType VarUint with wrong UintType for value
	ErrWrongUintType = errors.New("error var uint type with value")
)

const (
	// VarUint16 serialize type for uint16
	VarUint16 = 0xFD
	// VarUint32 serialize type for uint32
	VarUint32 = 0xFE
	// VarUint64 serialize type for uint64
	VarUint64 = 0xFF
)

// VarUint variable unsigned integer
type VarUint struct {
	UintType uint8 // VarUint8:[0, 0xFD); VarUint16: OxFD; VarUint32: 0xFE; VarUint64: 0xFF
	Value    uint64
}

// GetUintTypeByValue auxiliary function to get uint type by value
func GetUintTypeByValue(v uint64) uint8 {
	var t uint8 = 0x00
	if v < VarUint16 {
		t = uint8(v)
	} else if v >= VarUint16 && v <= math.MaxUint16 {
		t = VarUint16
	} else if v > math.MaxUint16 && v <= math.MaxUint32 {
		t = VarUint32
	} else /*if v > math.MaxUint32 && v <= math.MaxUint64 */ {
		t = VarUint64
	}
	return t
}

// Serialize implement Serializable interface
// serialze an variable unsigned integer with principle below:
// uint8:  serialize uint8 directly
// uin16:  0xFD + uint16(LE)
// uint32: 0xFE + uint32(LE)
// uint64: 0xFF + uint64(LE)
func (vu *VarUint) Serialize(w io.Writer) error {
	if vu.UintType < VarUint16 {
		binary.Write(w, binary.LittleEndian, uint8(vu.Value))
	} else {
		if vu.UintType == VarUint16 && vu.Value > VarUint16 && vu.Value <= math.MaxUint16 {
			binary.Write(w, binary.LittleEndian, vu.UintType)
			binary.Write(w, binary.LittleEndian, uint16(vu.Value))
		} else if vu.UintType == VarUint32 && vu.Value > math.MaxUint16 && vu.Value <= math.MaxUint32 {
			binary.Write(w, binary.LittleEndian, vu.UintType)
			binary.Write(w, binary.LittleEndian, uint32(vu.Value))
		} else if vu.UintType == VarUint64 && vu.Value > math.MaxUint32 && vu.Value <= math.MaxUint64 {
			binary.Write(w, binary.LittleEndian, vu.UintType)
			binary.Write(w, binary.LittleEndian, vu.Value)
		} else {
			return ErrWrongUintType
		}
	}
	return nil
}

// Deserialize implement Serializable interface
// deserialize variable unsigned integer
// see method Serialize as reference
func (vu *VarUint) Deserialize(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &vu.UintType)
	if vu.UintType == VarUint16 {
		var v uint16
		binary.Read(r, binary.LittleEndian, &v)
		vu.Value = uint64(v)
	} else if vu.UintType == VarUint32 {
		var v uint32
		binary.Read(r, binary.LittleEndian, &v)
		vu.Value = uint64(v)
	} else if vu.UintType == VarUint64 {
		var v uint64
		binary.Read(r, binary.LittleEndian, &v)
		vu.Value = uint64(v)
	} else {
		vu.Value = uint64(vu.UintType)
	}
	return nil
}
