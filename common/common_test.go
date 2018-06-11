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
	"reflect"
	"strings"
	"testing"
)

func TestNonce(t *testing.T) {
	r1 := Nonce()
	r2 := Nonce()
	if r1 == r2 {
		t.Errorf("nonce: %d, %d", r1, r2)
	}
}

func TestHex(t *testing.T) {
	b := []byte{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA}
	h := Hex(b)
	if strings.ToLower("FFFEFDFCFBFA") != h {
		t.Errorf("hex: %s", h)
	}
}

func TestHexToBytes(t *testing.T) {
	b := []byte{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA}
	h := Hex(b)
	bs, err := HexToBytes(h)
	if err != nil {
		t.Errorf("hex to bytes: %s", err)
	}
	if !reflect.DeepEqual(bs[:], b[:]) {
		t.Errorf("hex to bytes:\n0-%X,\n1-%X", b[:], bs[:])
	}
}

func TestFileExist(t *testing.T) {
	if FileExisted("dumy/foo.txt") {
		t.Errorf("common fileexits")
	}
}
