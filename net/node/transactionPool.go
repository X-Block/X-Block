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


func (txnPool *TXNPool) CleanTxnPool(txs []*transaction.Transaction) error {
	txsNum := len(txs)
	txInPoolNum := len(txnPool.list)
	cleaned := 0
	for _, tx := range txs[1:] {
		delete(txnPool.list, tx.Hash())
		cleaned++
	}
	if txsNum-cleaned != 1 {
		log.Info(fmt.Sprintf("The Transactions num Unmatched. Expect %d, got %d .\n", txsNum, cleaned))
	}
	log.Debug(fmt.Sprintf("[CleanTxnPool], Requested %d clean, %d transactions cleaned from localNode.TransPool and remains %d still in TxPool", txsNum, cleaned, txInPoolNum-cleaned))
	return nil
}

func (txnPool *TXNPool) init() {
	txnPool.list = make(map[common.Uint256]*transaction.Transaction)
	txnPool.txnCnt = 0
}

func (node *node) SynchronizeTxnPool() {
	node.nbrNodes.RLock()
	defer node.nbrNodes.RUnlock()

	for _, n := range node.nbrNodes.List {
		if n.state == ESTABLISH {
			msg.ReqTxnPool(n)
		}
	}
}

func (txnPool *TXNPool) CleanSubmittedTransactions(block *ledger.Block) error {
	txnPool.Lock()
	defer txnPool.Unlock()
	log.Debug()

	err := txnPool.CleanTxnPool(block.Transactions)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TxnPool], CleanSubmittedTransactions failed.")
	}
	return nil
}
