
package validate

import (
	"bytes"
	"encoding/binary"
	"io"
)

type mockVM struct {
	stack      []operand
	stackTop   int 
	origLength int 

	code *bytes.Reader

	polymorphic bool   
	blocks      []block 

	curFunc *wasm.FunctionSig
}

