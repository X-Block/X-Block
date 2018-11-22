package protocol

import (
	"XBlock/common"
	"XBlock/core/ledger"
	"XBlock/core/transaction"
	"XBlock/crypto"
	"XBlock/events"
	"bytes"
	"encoding/binary"
	"time"
)

type NodeAddr struct {
	Time     int64
	Services uint64
	IpAddr   [16]byte
	Port     uint16
	ID       uint64 
}


const (
	VERIFYNODE  = 1
	SERVICENODE = 2
)

const (
	VERIFYNODENAME  = "verify"
	SERVICENODENAME = "service"
)

const (
	MSGCMDLEN     = 12
	CMDOFFSET     = 4
	CHECKSUMLEN   = 4
	HASHLEN       = 32 
	MSGHDRLEN     = 24
	NETMAGIC      = 0x52597192
	MAXBLKHDRCNT  = 2000
	MAXINVHDRCNT  = 500
	DIVHASHLEN    = 5
	MINCONNCNT    = 3
	MAXREQBLKONCE = 16
)
const (
	HELLOTIMEOUT     = 3 
	MAXHELLORETYR    = 3
	MAXBUFLEN        = 1024 * 1024 * 5 
	MAXCHANBUF       = 512
	PROTOCOLVERSION  = 0
	PERIODUPDATETIME = 3 
	HEARTBEAT        = 2
	KEEPALIVETIMEOUT = 3
)


