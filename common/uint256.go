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
	"encoding/binary"
	"errors"
	"io"
)

/*
 * uint256: unsigned integer of 256 bit size
 * 32 * 8 = 256
 * 1 byte = 8 bit
 * actually this is base256
 */

// UINT256_SIZE 32 bytes
const UINT256_SIZE = 32

// Uint256 base256 with fix size byte array format
type Uint256 [UINT256_SIZE]byte

var (
	// UINT256_EMPTY empty uint256
	UINT256_EMPTY = Uint256{}

	// ErrBytesSize the byte array must be 32
	ErrBytesSize = errors.New("wrong bytes array size")
)

// Serialize implement Serializable interface
// serialize fix size byte array
func (u *Uint256) Serialize(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, u)
}

// Deserialize implement Serializable interface
// deserialize buffer to fix size byte array
func (u *Uint256) Deserialize(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, u)
}

// Bytes return bytes with copied content
func (u *Uint256) Bytes() []byte {
	b := make([]byte, UINT256_SIZE)
	copy(b, u[:])
	return b
}

// FromBytes copy content from bytes array to uint256
func (u *Uint256) FromBytes(b []byte) error {
	if len(b) != UINT256_SIZE {
		return ErrBytesSize
	}
	copy(u[:], b[:])
	return nil
}
