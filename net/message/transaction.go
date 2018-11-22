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

func (msg trn) Handle(node Noder) error {
	log.Debug()
	log.Debug("RX Transaction message")
	tx := &msg.txn
	if !node.LocalNode().ExistedID(tx.Hash()) {
		if err := va.VerifyTransaction(tx); err != nil {
			return errors.New("[VerifyTransaction] error")
		}
		if err := va.VerifyTransactionWithLedger(tx, ledger.DefaultLedger); err != nil {
			return errors.New("[VerifyTransactionWithLedger] error")
		}
		node.LocalNode().AppendTxnPool(&(msg.txn))
		node.LocalNode().IncRxTxnCnt()
		log.Debug("RX Transaction message hash", msg.txn.Hash())
		log.Debug("RX Transaction message type", msg.txn.TxType)
	}

	return nil
}

func reqTxnData(node Noder, hash common.Uint256) error {
	var msg dataReq
	msg.dataType = common.TRANSACTION
	
	buf, _ := msg.Serialization()
	go node.Tx(buf)
	return nil
}

func (msg dataReq) Serialization() ([]byte, error) {
	hdrBuf, err := msg.msgHdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	err = binary.Write(buf, binary.LittleEndian, msg.dataType)
	if err != nil {
		return nil, err
	}
	msg.hash.Serialize(buf)

	return buf.Bytes(), err
}

