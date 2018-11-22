package net

import (
	. "XBlock/common"
	"XBlock/core/ledger"
	"XBlock/core/transaction"
	"XBlock/crypto"
	"XBlock/events"
	"XBlock/net/node"
	"XBlock/net/protocol"
)

type Neter interface {
	GetTxnPool(cleanPool bool) map[Uint256]*transaction.Transaction
	SynchronizeTxnPool()
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookKeepersAddrs() ([]*crypto.PubKey, uint64)
	CleanSubmittedTransactions(block *ledger.Block) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
}

