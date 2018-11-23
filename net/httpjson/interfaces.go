package httpjsonrpc

import (
	"XBlock/account"
	. "XBlock/common"
	"XBlock/common/config"
	"XBlock/common/log"
	"XBlock/core/ledger"
	tx "XBlock/core/transaction"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
)

const (
	RANDBYTELEN = 4
)

func TransArryByteToHexString(ptx *tx.Transaction) *Transactions {

	trans := new(Transactions)
	trans.TxType = ptx.TxType
	trans.PayloadVersion = ptx.PayloadVersion
	trans.Payload = TransPayloadToHex(ptx.Payload)
	trans.Nonce = ptx.Nonce

	n := 0
	trans.Attributes = make([]TxAttributeInfo, len(ptx.Attributes))
	for _, v := range ptx.Attributes {
		trans.Attributes[n].Usage = v.Usage
		trans.Attributes[n].Data = ToHexString(v.Data)
		n++
	}

	n = 0
	trans.UTXOInputs = make([]UTXOTxInputInfo, len(ptx.UTXOInputs))
	for _, v := range ptx.UTXOInputs {
		trans.UTXOInputs[n].ReferTxID = ToHexString(v.ReferTxID.ToArray())
		trans.UTXOInputs[n].ReferTxOutputIndex = v.ReferTxOutputIndex
		n++
	}

	n = 0
	trans.BalanceInputs = make([]BalanceTxInputInfo, len(ptx.BalanceInputs))
	for _, v := range ptx.BalanceInputs {
		trans.BalanceInputs[n].AssetID = ToHexString(v.AssetID.ToArray())
		trans.BalanceInputs[n].Value = v.Value
		trans.BalanceInputs[n].ProgramHash = ToHexString(v.ProgramHash.ToArray())
		n++
	}

	n = 0
	trans.Outputs = make([]TxoutputInfo, len(ptx.Outputs))
	for _, v := range ptx.Outputs {
		trans.Outputs[n].AssetID = ToHexString(v.AssetID.ToArray())
		trans.Outputs[n].Value = v.Value
		trans.Outputs[n].ProgramHash = ToHexString(v.ProgramHash.ToArray())
		n++
	}

	n = 0
	trans.Programs = make([]ProgramInfo, len(ptx.Programs))
	for _, v := range ptx.Programs {
		trans.Programs[n].Code = ToHexString(v.Code)
		trans.Programs[n].Parameter = ToHexString(v.Parameter)
		n++
	}

	n = 0
	trans.AssetOutputs = make([]TxoutputMap, len(ptx.AssetOutputs))
	for k, v := range ptx.AssetOutputs {
		trans.AssetOutputs[n].Key = k
		trans.AssetOutputs[n].Txout = make([]TxoutputInfo, len(v))
		for m := 0; m < len(v); m++ {
			trans.AssetOutputs[n].Txout[m].AssetID = ToHexString(v[m].AssetID.ToArray())
			trans.AssetOutputs[n].Txout[m].Value = v[m].Value
			trans.AssetOutputs[n].Txout[m].ProgramHash = ToHexString(v[m].ProgramHash.ToArray())
		}
		n += 1
	}

	n = 0
	trans.AssetInputAmount = make([]AmountMap, len(ptx.AssetInputAmount))
	for k, v := range ptx.AssetInputAmount {
		trans.AssetInputAmount[n].Key = k
		trans.AssetInputAmount[n].Value = v
		n += 1
	}

	n = 0
	trans.AssetOutputAmount = make([]AmountMap, len(ptx.AssetOutputAmount))
	for k, v := range ptx.AssetOutputAmount {
		trans.AssetInputAmount[n].Key = k
		trans.AssetInputAmount[n].Value = v
		n += 1
	}

	mhash := ptx.Hash()
	trans.Hash = ToHexString(mhash.ToArray())

	return trans
}

func getBestBlockHash(params []interface{}) map[string]interface{} {
	hash := ledger.DefaultLedger.Blockchain.CurrentBlockHash()
	return XBlockRpc(ToHexString(hash.ToArray()))
}

func getBlock(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return XBlockRpcNil
	}
	var err error
	var hash Uint256
	switch (params[0]).(type) {
	case float64:
		index := uint32(params[0].(float64))
		hash, err = ledger.DefaultLedger.Store.GetBlockHash(index)
		if err != nil {
			return XBlockRpcUnknownBlock
		}
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return XBlockRpcInvalidParameter
		}
		if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
			return XBlockRpcInvalidTransaction
		}
	default:
		return XBlockRpcInvalidParameter
	}

	block, err := ledger.DefaultLedger.Store.GetBlock(hash)
	if err != nil {
		return XBlockRpcUnknownBlock
	}

	blockHead := &BlockHead{
		Version:          block.Blockdata.Version,
		PrevBlockHash:    ToHexString(block.Blockdata.PrevBlockHash.ToArray()),
		TransactionsRoot: ToHexString(block.Blockdata.TransactionsRoot.ToArray()),
		Timestamp:        block.Blockdata.Timestamp,
		Height:           block.Blockdata.Height,
		ConsensusData:    block.Blockdata.ConsensusData,
		NextBookKeeper:   ToHexString(block.Blockdata.NextBookKeeper.ToArray()),
		Program: ProgramInfo{
			Code:      ToHexString(block.Blockdata.Program.Code),
			Parameter: ToHexString(block.Blockdata.Program.Parameter),
		},
		Hash: ToHexString(hash.ToArray()),
	}

	trans := make([]*Transactions, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		trans[i] = TransArryByteToHexString(block.Transactions[i])
	}

	b := BlockInfo{
		Hash:         ToHexString(hash.ToArray()),
		BlockData:    blockHead,
		Transactions: trans,
	}
	return XBlockRpc(b)
}

func getBlockCount(params []interface{}) map[string]interface{} {
	return XBlockRpc(ledger.DefaultLedger.Blockchain.BlockHeight + 1)
}

func getBlockHash(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return XBlockRpcNil
	}
	switch params[0].(type) {
	case float64:
		height := uint32(params[0].(float64))
		hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
		if err != nil {
			return XBlockRpcUnknownBlock
		}
		return XBlockRpc(fmt.Sprintf("%016x", hash))
	default:
		return XBlockRpcInvalidParameter
	}
}

func getConnectionCount(params []interface{}) map[string]interface{} {
	return XBlockRpc(node.GetConnectionCnt())
}

func getRawMemPool(params []interface{}) map[string]interface{} {
	txs := []*Transactions{}
	txpool := node.GetTxnPool(false)
	for _, t := range txpool {
		txs = append(txs, TransArryByteToHexString(t))
	}
	if len(txs) == 0 {
		return XBlockRpcNil
	}
	return XBlockRpc(txs)
}

func getRawTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return XBlockRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return XBlockRpcInvalidParameter
		}
		var hash Uint256
		err = hash.Deserialize(bytes.NewReader(hex))
		if err != nil {
			return XBlockRpcInvalidTransaction
		}
		tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
		if err != nil {
			return XBlockRpcUnknownTransaction
		}
		tran := TransArryByteToHexString(tx)
		return XBlockRpc(tran)
	default:
		return XBlockRpcInvalidParameter
	}
}

func sendRawTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return XBlockRpcNil
	}
	var hash Uint256
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return XBlockRpcInvalidParameter
		}
		var txn tx.Transaction
		if err := txn.Deserialize(bytes.NewReader(hex)); err != nil {
			return XBlockRpcInvalidTransaction
		}
		hash = txn.Hash()
		if err := VerifyAndSendTx(&txn); err != nil {
			return XBlockRpcInternalError
		}
	default:
		return XBlockRpcInvalidParameter
	}
	return XBlockRpc(ToHexString(hash.ToArray()))
}

func getUnspendOutput(params []interface{}) map[string]interface{} {
	if len(params) < 2 {
		return XBlockRpcNil
	}
	var programHash Uint160
	var assetHash Uint256
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return XBlockRpcInvalidParameter
		}
		if err := programHash.Deserialize(bytes.NewReader(hex)); err != nil {
			return XBlockRpcInvalidHash
		}
	default:
		return XBlockRpcInvalidParameter
	}

	switch params[1].(type) {
	case string:
		str := params[1].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return XBlockRpcInvalidParameter
		}
		if err := assetHash.Deserialize(bytes.NewReader(hex)); err != nil {
			return XBlockRpcInvalidHash
		}
	default:
		return XBlockRpcInvalidParameter
	}
	type TxOutputInfo struct {
		AssetID     string
		Value       int64
		ProgramHash string
	}
	outputs := make(map[string]*TxOutputInfo)
	height := ledger.DefaultLedger.GetLocalBlockChainHeight()
	var i uint32
	for i = 0; i <= height; i++ {
		block, err := ledger.DefaultLedger.GetBlockWithHeight(i)
		if err != nil {
			return XBlockRpcInternalError
		}
		for _, t := range block.Transactions[1:] {
			if t.TxType == tx.RegisterAsset {
				continue
			}
			txHash := t.Hash()
			txHashHex := ToHexString(txHash.ToArray())
			for i, output := range t.Outputs {
				if output.AssetID.CompareTo(assetHash) == 0 &&
					output.ProgramHash.CompareTo(programHash) == 0 {
					key := txHashHex + ":" + strconv.Itoa(i)
					asset := ToHexString(output.AssetID.ToArray())
					pHash := ToHexString(output.ProgramHash.ToArray())
					value := int64(output.Value)
					info := &TxOutputInfo{
						asset,
						value,
						pHash,
					}
					outputs[key] = info
				}
			}
		}
	}
	height = ledger.DefaultLedger.GetLocalBlockChainHeight()
	for i = 0; i <= height; i++ {
		block, err := ledger.DefaultLedger.GetBlockWithHeight(i)
		if err != nil {
			return XBlockRpcInternalError
		}
		for _, t := range block.Transactions[1:] {
			if t.TxType == tx.RegisterAsset {
				continue
			}
			for _, input := range t.UTXOInputs {
				refer := ToHexString(input.ReferTxID.ToArray())
				index := strconv.Itoa(int(input.ReferTxOutputIndex))
				key := refer + ":" + index
				delete(outputs, key)
			}
		}
	}
	return XBlockRpc(outputs)
}

func getTxout(params []interface{}) map[string]interface{} {
	return XBlockRpcUnsupported
}

func submitBlock(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return XBlockRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, _ := hex.DecodeString(str)
		var block ledger.Block
		if err := block.Deserialize(bytes.NewReader(hex)); err != nil {
			return XBlockRpcInvalidBlock
		}
		if err := ledger.DefaultLedger.Blockchain.AddBlock(&block); err != nil {
			return XBlockRpcInvalidBlock
		}
		if err := node.Xmit(&block); err != nil {
			return XBlockRpcInternalError
		}
	default:
		return XBlockRpcInvalidParameter
	}
	return XBlockRpcSuccess
}

func getNeighbor(params []interface{}) map[string]interface{} {
	addr, _ := node.GetNeighborAddrs()
	return XBlockRpc(addr)
}

func getNodeState(params []interface{}) map[string]interface{} {
	n := NodeInfo{
		State:    uint(node.GetState()),
		Time:     node.GetTime(),
		Port:     node.GetPort(),
		ID:       node.GetID(),
		Version:  node.Version(),
		Services: node.Services(),
		Relay:    node.GetRelay(),
		Height:   node.GetHeight(),
		TxnCnt:   node.GetTxnCnt(),
		RxTxnCnt: node.GetRxTxnCnt(),
	}
	return XBlockRpc(n)
}

func startConsensus(params []interface{}) map[string]interface{} {
	if err := dBFT.Start(); err != nil {
		return XBlockRpcFailed
	}
	return XBlockRpcSuccess
}

func stopConsensus(params []interface{}) map[string]interface{} {
	if err := dBFT.Halt(); err != nil {
		return XBlockRpcFailed
	}
	return XBlockRpcSuccess
}

func sendSampleTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return XBlockRpcNil
	}
	var txType string
	switch params[0].(type) {
	case string:
		txType = params[0].(string)
	default:
		return XBlockRpcInvalidParameter
	}

	issuer, err := account.NewAccount()
	if err != nil {
		return XBlockRpc("Failed to create account")
	}
	admin := issuer

	rbuf := make([]byte, RANDBYTELEN)
	rand.Read(rbuf)
	switch string(txType) {
	case "perf":
		num := 1
		if len(params) == 2 {
			switch params[1].(type) {
			case float64:
				num = int(params[1].(float64))
			}
		}
		for i := 0; i < num; i++ {
			regTx := NewRegTx(ToHexString(rbuf), i, admin, issuer)
			SignTx(admin, regTx)
			VerifyAndSendTx(regTx)
		}
		return XBlockRpc(fmt.Sprintf("%d transaction(s) was sent", num))
	case "bookkeeper":

		if len(params) < 3 {
			return XBlockRpcNil
		}
		ind := 4
		switch params[1].(type) {
		case float64:
			ind = int(params[1].(float64))
		}
		var isAdd bool
		switch params[2].(type) {
		case string:
			action := params[2].(string)
			if action == "add" {
				isAdd = true
			} else if action == "sub" {
				isAdd = false
			} else {
				return XBlockRpcInvalidParameter
			}
		default:
			return XBlockRpcInvalidParameter
		}

		walletFile := "wallet" + strconv.Itoa(ind) + ".dat"
		c := account.Open(walletFile, []byte(account.DefaultPin))
		if c == nil {
			return XBlockRpc("do not have wallet file:" + walletFile)
		}

		account, _ := c.GetDefaultAccount()
		pubKey := account.PubKey()

		cert := make([]byte, 100)
		rand.Read(cert)

		bkTx, _ := tx.NewBookKeeperTransaction(pubKey, isAdd, cert)
		VerifyAndSendTx(bkTx)

		return XBlockRpc(fmt.Sprint("bookkeeper transaction was sent, select pubkey file:", walletFile))

	default:
		return XBlockRpc("Invalid transacion type")
	}
}

