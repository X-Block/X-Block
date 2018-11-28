
package wasm

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/X-Block/X-Block/tree/master/vm/wasm/leb128"
)

const currentVersion = 0x01


func EncodeModule(w io.Writer, m *Module) error {
	if err := writeU32(w, Magic); err != nil {
		return err
	}
	if err := writeU32(w, currentVersion); err != nil {
		return err
	}
	sections := m.Sections
	buf := new(bytes.Buffer)
	for _, s := range sections {
		if _, err := leb128.WriteVarUint32(w, uint32(s.SectionID())); err != nil {
			return err
		}
		buf.Reset()
		if err := s.WritePayload(buf); err != nil {
			return err
		}
		if _, err := leb128.WriteVarUint32(w, uint32(buf.Len())); err != nil {
			return err
		}
		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
	}
	return nil
}

func writeStringUint(w io.Writer, s string) error {
	return writeBytesUint(w, []byte(s))
}

