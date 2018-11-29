
package wasm

import (
	"errors"
	"io"
	"reflect"

	"github.com/X-Block/X-Block/tree/master/vm/wasm/internal/readpos"
)

var ErrInvalidMagic = errors.New("wasm: Invalid magic number")

type Function struct {
	Sig  *FunctionSig
	Body *FunctionBody
	Host reflect.Value
}

