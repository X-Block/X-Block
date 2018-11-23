package message

import (
	"XBlock/common"
	"XBlock/common/log"
	"XBlock/core/ledger"
	"XBlock/events"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

type blockReq struct {
	msgHdr
}

type block struct {
	msgHdr
	blk ledger.Block
}

func (msg block) Handle(node Noder) error {
	log.Debug("RX block message")
	hash := msg.blk.Hash()
	if ledger.DefaultLedger.BlockInLedger(hash) {
		log.Warn("Receive duplicated block: ", hash)
		return errors.New("Received duplicate block")
	}
	if err := ledger.DefaultLedger.Blockchain.AddBlock(&msg.blk); err != nil {
		log.Error("Block adding error: ", hash)
		return err
	}
	node.RemoveFlightHeight(msg.blk.Blockdata.Height)
	node.LocalNode().GetEvent("block").Notify(events.EventNewInventory, &msg.blk)
	return nil
}

func (msg dataReq) Handle(node Noder) error {
	log.Debug()
	reqtype := common.InventoryType(msg.dataType)
	hash := msg.hash
	switch reqtype {
	case common.BLOCK:
		block, err := NewBlockFromHash(hash)
		if err != nil {
			log.Error("Can't get block from hash: ", hash, " ,send not found message")
			b, err := NewNotFound(hash)
			node.Tx(b)
			return err
		}
		log.Debug("block height is ", block.Blockdata.Height, " ,hash is ", hash)
		buf, err := NewBlock(block)
		if err != nil {
			return err
		}
		node.Tx(buf)

	case common.TRANSACTION:
		txn, err := NewTxnFromHash(hash)
		if err != nil {
			return err
		}
		buf, err := NewTxn(txn)
		if err != nil {
			return err
		}
		go node.Tx(buf)
	}
	return nil
}

