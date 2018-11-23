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

