
package wasm

import (
	"io/ioutil"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	SetDebugMode(false)
}

func SetDebugMode(dbg bool) {
	w := ioutil.Discard
	if dbg {
		w = os.Stderr
	}
	logger = log.New(w, "", log.Lshortfile)
}
