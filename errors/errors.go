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

package errors

import "errors"

const (
	callStackDepth = 10
)

//DetailError error with details
type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetErrRoot() error
}

//NewErr create a new error with error message
func NewErr(errmsg string) error {
	return errors.New(errmsg)
}

//NewDetailErr create new DetailError
func NewDetailErr(err error, errcode ErrCode, errmsg string) DetailError {
	if err == nil {
		return nil
	}

	e, ok := err.(Errors)
	if !ok {
		e.root = err
		e.errmsg = err.Error()
		e.callstack = getCallStack(0, callStackDepth)
		e.code = errcode

	}
	if errmsg != "" {
		e.errmsg = errmsg + ": " + e.errmsg
	}

	return e
}

//RootErr error root
func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetErrRoot()
	}
	return err
}

// Errors errors
type Errors struct {
	errmsg    string
	callstack *CallStack
	root      error
	code      ErrCode
}

//Error get error message
func (e Errors) Error() string {
	return e.errmsg
}

//GetErrCode get error code
func (e Errors) GetErrCode() ErrCode {
	return e.code
}

//GetErrRoot implement DetailError interface
func (e Errors) GetErrRoot() error {
	return e.root
}

//GetCallStack get call stack
func (e Errors) GetCallStack() *CallStack {
	return e.callstack
}
