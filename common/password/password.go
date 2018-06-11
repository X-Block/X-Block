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

package password

import (
	"bytes"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

// WaitForPwdInput wait for user to input password
func WaitForPwdInput() ([]byte, error) {
	fmt.Printf("Please input your Password:")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return passwd, nil
}

// PwdInputAndConfirm wait for user to input password and its confirmation
func PwdInputAndConfirm() ([]byte, error) {
	fmt.Printf("Please input your Password:")
	first, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if 0 == len(first) {
		fmt.Println("You have to input password.")
		os.Exit(1)
	}

	fmt.Printf("Re-enter Password:")
	second, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if 0 == len(second) {
		fmt.Println("You have to input password.")
		os.Exit(1)
	}

	if !bytes.Equal(first, second) {
		fmt.Println("Unmatched Password")
		os.Exit(1)
	}
	return first, nil
}
