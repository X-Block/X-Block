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
	"encoding/hex"
	"math/rand"
	"os"
)

// Nonce returns random nonce
func Nonce() uint64 {
	// TODO: replace with the real random number generator
	nonce := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	return nonce
}

// Hex get hex string corresponding to byte array
func Hex(data []byte) string {
	return hex.EncodeToString(data)
}

// HexToBytes convert hex string to byte array
func HexToBytes(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

// FileExisted checks whether filename exists in filesystem
func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
