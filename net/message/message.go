package message

import (
	"XBlock/common/log"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type Messager interface {
	Verify([]byte) error
	Serialization() ([]byte, error)
	Deserialization([]byte) error
	Handle(Noder) error
}

type msgHdr struct {
	Magic uint32
	
	CMD      [MSGCMDLEN]byte 
	Length   uint32
	Checksum [CHECKSUMLEN]byte
}


type msgCont struct {
	hdr msgHdr
	p   interface{}
}

type varStr struct {
	len uint
	buf []byte
}

type filteradd struct {
	msgHdr
	
}

type filterclear struct {
	msgHdr
	
}

type filterload struct {
	msgHdr
	
}

