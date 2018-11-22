package message

import (
	"XBlock/common"
	"XBlock/common/log"
	"XBlock/core/ledger"
	"XBlock/core/transaction"
	va "XBlock/core/validation"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

type dataReq struct {
	msgHdr
	dataType common.InventoryType
	hash     common.Uint256
}


type trn struct {
	msgHdr
	txn transaction.Transaction
}

