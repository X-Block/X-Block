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

