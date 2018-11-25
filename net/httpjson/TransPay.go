package httpjsonrpc

import (
	. "XBlock/common"
	"XBlock/core/asset"
	. "XBlock/core/contract"
	. "XBlock/core/transaction"
	"XBlock/core/transaction/payload"
	"bytes"
)

type Payload interface {
	Data() []byte
}


type BookKeeping struct {
}

func (dc *BookKeeping) Data() []byte {
	return []byte{0}
}


