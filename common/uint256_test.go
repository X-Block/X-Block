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

package common

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUint256Serialize(t *testing.T) {
	b := new(bytes.Buffer)
	var u256 = Uint256{0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD}
	if err := u256.Serialize(b); err != nil {
		t.Errorf("uint256 seriliaze: %s", err)
	}
	bs := b.Bytes()
	if !reflect.DeepEqual(u256[:], bs[:]) {
		t.Errorf("uint256 serialize: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}

func TestUint256Deserialize(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD})
	var u256 = UINT256_EMPTY
	bs := b.Bytes()
	if err := u256.Deserialize(b); err != nil {
		t.Errorf("uint256 deserialize: %s", err)
	}

	if !reflect.DeepEqual(u256[:], bs[:]) {
		t.Errorf("uint256 serialize: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}

func TestUint256Bytes(t *testing.T) {
	var u256 = Uint256{0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD}
	bs := u256.Bytes()
	if !reflect.DeepEqual(bs[:], u256[:]) {
		t.Errorf("uint256 get bytes: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}

func TestUint256FromBytes(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD})
	var u256 = UINT256_EMPTY
	bs := b.Bytes()
	if err := u256.FromBytes(bs); err != nil {
		t.Errorf("uint256 frombytes: %s", err)
	}

	if !reflect.DeepEqual(u256[:], bs[:]) {
		t.Errorf("uint256 frombytes: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}
