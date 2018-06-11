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

func TestVarBytesSerialize(t *testing.T) {
	b := new(bytes.Buffer)
	bts := make([]byte, 4)
	bts[0] = 't'
	bts[1] = 'e'
	bts[2] = 's'
	bts[3] = 't'
	var varbt = VarBytes{Len: uint64(len(bts)), Bytes: bts}
	if err := varbt.Serialize(b); err != nil {
		t.Errorf("varbytes: %s", err)
	}
	bs := b.Bytes()
	if uint8(bs[0]) != 0x04 ||
		bs[1] != 't' ||
		bs[2] != 'e' ||
		bs[3] != 's' ||
		bs[4] != 't' {
		t.Errorf("varbytes: 0-%X, 1-%c, 2-%c, 3-%c, 4-%c", uint8(bs[0]), bs[1], bs[2], bs[3], bs[4])
	}
}

func TestVarBytesDeserialize(t *testing.T) {
	b := new(bytes.Buffer)
	bts := make([]byte, 4)
	bts[0] = 't'
	bts[1] = 'e'
	bts[2] = 's'
	bts[3] = 't'
	var varbt = VarBytes{Len: uint64(len(bts)), Bytes: bts}
	if err := varbt.Serialize(b); err != nil {
		t.Errorf("varbytes: %s", err)
	}

	var vardbt VarBytes
	if err := vardbt.Deserialize(b); err != nil {
		t.Errorf("varbytes; %s", err)
	}
	if vardbt.Len != 0x04 ||
		vardbt.Bytes[0] != 't' ||
		vardbt.Bytes[1] != 'e' ||
		vardbt.Bytes[2] != 's' ||
		vardbt.Bytes[3] != 't' {
		t.Errorf("varbytes: 0-%X, 1-%c, 2-%c, 3-%c, 4-%c", vardbt.Len, vardbt.Bytes[0], vardbt.Bytes[1], vardbt.Bytes[2], vardbt.Bytes[3])
	}

}
