
package main

import (
	"bytes"
	"encoding/hex"
	"io"
)


func hexDump(data []byte, offset uint) string {
	buf := new(bytes.Buffer)
	d := &dumper{w: buf, n: offset}
	d.Write(data)
	d.Close()
	return buf.String()
}

type dumper struct {
	w          io.Writer
	rightChars [18]byte
	buf        [14]byte
	used       int  
	n          uint 
}

