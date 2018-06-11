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
	"bytes"
	"testing"
)

func TestVarUintSerialize(t *testing.T) {
	b := new(bytes.Buffer)
	a8 := uint8(0xF0)
	a16 := uint16(0xF0FE)
	a32 := uint32(0xFFF0FFFE)
	a64 := uint64(0xFFF0FFF1FFF2FFFE)

	var varuint8 = VarUint{UintType: GetUintTypeByValue(uint64(a8)), Value: uint64(a8)}
	if err := varuint8.Serialize(b); err != nil {
		t.Errorf("varuint8: %s", err)
	} else {
		t.Log("varuint8 passed")
	}
	var varuint16 = VarUint{UintType: GetUintTypeByValue(uint64(a16)), Value: uint64(a16)}
	if err := varuint16.Serialize(b); err != nil {
		t.Errorf("varuint16: %s", err)
	} else {
		t.Log("varuint16 passed")
	}
	var varuint32 = VarUint{UintType: GetUintTypeByValue(uint64(a32)), Value: uint64(a32)}
	if err := varuint32.Serialize(b); err != nil {
		t.Errorf("varuint32: %s", err)
	} else {
		t.Log("varuint32 passed")
	}
	var varuint64 = VarUint{UintType: GetUintTypeByValue(uint64(a64)), Value: uint64(a64)}
	if err := varuint64.Serialize(b); err != nil {
		t.Errorf("varuint64: %s", err)
	} else {
		t.Log("varuint64 passed")
	}
	ba := b.Bytes()
	if uint8(ba[0]) != 0xF0 {
		t.Errorf("varuint8: %X", uint8(ba[0]))
	}
	if uint8(ba[1]) != VarUint16 &&
		uint8(ba[2]) != 0xFE &&
		uint8(ba[3]) != 0xF0 {
		t.Errorf("varuint16: 0-%X, 1-%X, 2-%X",
			uint8(ba[1]), uint8(ba[2]), uint8(ba[3]))
	}
	if uint8(ba[4]) != VarUint32 &&
		uint8(ba[5]) != 0xFE &&
		uint8(ba[6]) != 0xFF &&
		uint8(ba[7]) != 0xF0 &&
		uint8(ba[8]) != 0xFF {
		t.Errorf("varuint32: 0-%X, 1-%X, 2-%X, 3-%X, 4-%X",
			uint8(ba[4]), uint8(ba[5]), uint8(ba[6]), uint8(ba[7]), uint8(ba[8]))
	}
	if uint8(ba[9]) != VarUint64 &&
		uint8(ba[10]) != 0xFE &&
		uint8(ba[11]) != 0xFF &&
		uint8(ba[12]) != 0xF2 &&
		uint8(ba[13]) != 0xFF &&
		uint8(ba[14]) != 0xF1 &&
		uint8(ba[15]) != 0xFF &&
		uint8(ba[16]) != 0xF0 &&
		uint8(ba[17]) != 0xFF {
		t.Errorf("varuint64: 0-%X, 1-%X, 2-%X, 3-%X, 4-%X, 5-%X, 6-%X, 7-%X, 8-%X",
			uint8(ba[9]), uint8(ba[10]), uint8(ba[11]), uint8(ba[12]), uint8(ba[13]), uint8(ba[14]), uint8(ba[15]), uint8(ba[16]), uint8(ba[17]))
	}

}

func TestVarUintDerialize(t *testing.T) {
	b := new(bytes.Buffer)
	a8 := uint8(0xF0)
	a16 := uint16(0xF0FE)
	a32 := uint32(0xFFF0FFFE)
	a64 := uint64(0xFFF0FFF1FFF2FFFE)
	var varuint8 = VarUint{UintType: a8, Value: uint64(a8)}
	if err := varuint8.Serialize(b); err != nil {
		t.Errorf("varuint8: %s", err)
	} else {
		t.Log("varuint8 passed")
	}
	var varuint16 = VarUint{UintType: VarUint16, Value: uint64(a16)}
	if err := varuint16.Serialize(b); err != nil {
		t.Errorf("varuint16: %s", err)
	} else {
		t.Log("varuint16 passed")
	}
	var varuint32 = VarUint{UintType: VarUint32, Value: uint64(a32)}
	if err := varuint32.Serialize(b); err != nil {
		t.Errorf("varuint32: %s", err)
	} else {
		t.Log("varuint32 passed")
	}
	var varuint64 = VarUint{UintType: VarUint64, Value: uint64(a64)}
	if err := varuint64.Serialize(b); err != nil {
		t.Errorf("varuint64: %s", err)
	} else {
		t.Log("varuint64 passed")
	}

	var dvu8 VarUint
	if err := dvu8.Deserialize(b); err != nil {
		t.Errorf("deserialize varuint8: %s", err)
	}
	if dvu8.UintType != 0xF0 || dvu8.Value != 0xF0 {
		t.Errorf("deserialize varuint8: 0-%X, 1-%X", dvu8.UintType, dvu8.Value)
	}
	var dvu16 VarUint
	if err := dvu16.Deserialize(b); err != nil {
		t.Errorf("deserialize varuint16: %s", err)
	}
	if dvu16.UintType != VarUint16 || dvu16.Value != 0xF0FE {
		t.Errorf("deserialize varuint16: 0-%X, 1-%X", dvu16.UintType, dvu16.Value)
	}
	var dvu32 VarUint
	if err := dvu32.Deserialize(b); err != nil {
		t.Errorf("deserialize varuint32: %s", err)
	}
	if dvu32.UintType != VarUint32 || dvu32.Value != 0xFFF0FFFE {
		t.Errorf("deserialize varuint32: 0-%X, 1-%X", dvu32.UintType, dvu32.Value)
	}
	var dvu64 VarUint
	if err := dvu64.Deserialize(b); err != nil {
		t.Errorf("deserialize varuint64: %s", err)
	}
	if dvu64.UintType != VarUint64 || dvu64.Value != 0xFFF0FFF1FFF2FFFE {
		t.Errorf("deserialize varuint64: 0-%X, 1-%X", dvu64.UintType, dvu64.Value)
	}

}
