package message

import (
	"XBlock/common/log"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
)

type addrReq struct {
	Hdr msgHdr
	
}

type addr struct {
	hdr       msgHdr
	nodeCnt   uint64
	nodeAddrs []NodeAddr
}

const (
	NODEADDRSIZE = 30
)

func newGetAddr() ([]byte, error) {
	var msg addrReq
	var sum []byte
	sum = []byte{0x23, 0x72, 0x35, 0x24}
	msg.Hdr.init("getaddr", sum, 0)

	buf, err := msg.Serialization()
	if err != nil {
		return nil, err
	}

	str := hex.EncodeToString(buf)
	log.Debug("The message get addr length is: ", len(buf), " ", str)

	return buf, err
}

