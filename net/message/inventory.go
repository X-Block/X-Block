package message

import (
	. "XBlock/common"
	"XBlock/common/log"
	"XBlock/common/serialization"
	"XBlock/core/ledger"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

type blocksReq struct {
	msgHdr
	p struct {
		HeaderHashCount uint8
		hashStart       [HASHLEN]byte
		hashStop        [HASHLEN]byte
	}
}

