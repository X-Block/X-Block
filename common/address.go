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
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"

	base58 "github.com/itchyny/base58-go"
)

const (
	// ADDR_LEN address length
	ADDR_LEN = 20
)

var (
	// ADDRESS_EMPTY empty address
	ADDRESS_EMPTY = Address{}

	// ErrBase58Addr wrong base58 address
	ErrBase58Addr = errors.New("wrong encoded address")
	// ErrBase58Verify verify base58 address error
	ErrBase58Verify = errors.New("base58 address verify")
)

// Address 20 byte length array
type Address [ADDR_LEN]byte

// Serialize implement Serializable interface
// serialize address byte array
func (addr *Address) Serialize(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, addr[:])
}

// Deserialize implement Serializable interface
// deserialize address byte array
func (addr *Address) Deserialize(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, addr)
}

// Hex get hex string corresponding to the Address
func (addr *Address) Hex() string {
	return fmt.Sprintf("%x", addr[:])
}

// Base58 get base58 string corresponding to the Address
// with principle below:
// data = 0x41 + address[:] + sha256(sha256(0x41 + address[:]))[0:4]
// base58 = encode(data)
func (addr *Address) Base58() string {
	data := append([]byte{0x41}, addr[:]...)
	temp := sha256.Sum256(data)
	temps := sha256.Sum256(temp[:])
	data = append(data, temps[0:4]...)

	bi := new(big.Int).SetBytes(data).String()
	encoded, _ := base58.BitcoinEncoding.Encode([]byte(bi))
	return string(encoded)
}

// FromBytes get Address from byte array
func (addr *Address) FromBytes(b []byte) error {
	if len(b) != ADDR_LEN {
		return ErrBytesSize
	}
	copy(addr[:], b)
	return nil
}

// FromBase58 get Address from Base58 address
// see 'Base58' method above as reference
func (addr *Address) FromBase58(base58addr string) error {
	decoded, err := base58.BitcoinEncoding.Decode([]byte(base58addr))
	if err != nil {
		return err
	}

	x, _ := new(big.Int).SetString(string(decoded), 10)

	buf := x.Bytes()
	if len(buf) != 1+ADDR_LEN+4 || buf[0] != byte(0x41) {
		return ErrBase58Addr
	}
	err = addr.FromBytes(buf[1:21])
	if err != nil {
		return err
	}

	addr58 := addr.Base58()

	if addr58 != base58addr {
		return ErrBase58Verify
	}

	return nil
}
