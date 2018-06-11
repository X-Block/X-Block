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
	"io"
)

// VarBytes variable bytes array with length
// serialized with format below:
// length + bytes[]
type VarBytes struct {
	Len   uint64
	Bytes []byte
}

// Serialize implement Serializable interface
// serialize a variable bytes array into format below:
// variable length + bytes[]
func (vb *VarBytes) Serialize(w io.Writer) error {
	var varlen = VarUint{UintType: GetUintTypeByValue(vb.Len), Value: vb.Len}
	if err := varlen.Serialize(w); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, vb.Bytes)
}

// Deserialize implement Deserialiazable interface
// deserialize a variable bytes arrary from buffer
// see Serialize above as reference
func (vb *VarBytes) Deserialize(r io.Reader) error {
	var varlen VarUint
	if err := varlen.Deserialize(r); err != nil {
		return err
	}
	vb.Len = uint64(varlen.Value)
	vb.Bytes = make([]byte, vb.Len)
	return binary.Read(r, binary.LittleEndian, vb.Bytes)
}
