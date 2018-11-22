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


