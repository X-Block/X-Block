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

func (txnPool *TXNPool) AppendTxnPool(txn *transaction.Transaction) bool {
	hash := txn.Hash()
	txnPool.Lock()
	txnPool.list[hash] = txn
	txnPool.txnCnt++
	txnPool.Unlock()
	return true
}


func (txnPool *TXNPool) GetTxnPool(cleanPool bool) map[common.Uint256]*transaction.Transaction {
	txnPool.Lock()
	defer txnPool.Unlock()

	list := txnPool.list
	if cleanPool == true {
		txnPool.init()
	}
	return DeepCopy(list)
}

func DeepCopy(mapIn map[common.Uint256]*transaction.Transaction) map[common.Uint256]*transaction.Transaction {
	reply := make(map[common.Uint256]*transaction.Transaction)
	for k, v := range mapIn {
		reply[k] = v
	}
	return reply
}


