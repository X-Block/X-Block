package node

import (
	"XBlock/common"
	"XBlock/common/log"
	"XBlock/core/ledger"
	"XBlock/core/transaction"
	"XBlock/errors"
	msg "XBlock/net/message"
	. "XBlock/net/protocol"
	"fmt"
	"sync"
)

type TXNPool struct {
	sync.RWMutex
	txnCnt uint64
	list   map[common.Uint256]*transaction.Transaction
}

func (txnPool *TXNPool) GetTransaction(hash common.Uint256) *transaction.Transaction {
	txnPool.RLock()
	defer txnPool.RUnlock()
	txn := txnPool.list[hash]
	return txn
}

