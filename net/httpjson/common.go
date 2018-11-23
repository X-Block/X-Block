package httpjson

import (
	. "XBlock/common"
	"XBlock/common/log"
	"XBlock/consensus/dbft"
	"XBlock/core/ledger"
	. "XBlock/core/transaction"
	tx "XBlock/core/transaction"
	"XBlock/core/validation"
	. "XBlock/net/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func init() {
	mainMux.m = make(map[string]func([]interface{}) map[string]interface{})
}


var mainMux ServeMux
var node Noder
var dBFT *dbft.DbftService

rpc call
type ServeMux struct {
	sync.RWMutex
	m               map[string]func([]interface{}) map[string]interface{}
	defaultFunction func(http.ResponseWriter, *http.Request)
}

type TxAttributeInfo struct {
	Usage TransactionAttributeUsage
	Data  string
}

type UTXOTxInputInfo struct {
	ReferTxID          string
	ReferTxOutputIndex uint16
}

type BalanceTxInputInfo struct {
	AssetID     string
	Value       Fixed64
	ProgramHash string
}

type TxoutputInfo struct {
	AssetID     string
	Value       Fixed64
	ProgramHash string
}

type TxoutputMap struct {
	Key   Uint256
	Txout []TxoutputInfo
}

type AmountMap struct {
	Key   Uint256
	Value Fixed64
}

type ProgramInfo struct {
	Code      string
	Parameter string
}

type Transactions struct {
	TxType         TransactionType
	PayloadVersion byte
	Payload        PayloadInfo
	Nonce          uint64
	Attributes     []TxAttributeInfo
	UTXOInputs     []UTXOTxInputInfo
	BalanceInputs  []BalanceTxInputInfo
	Outputs        []TxoutputInfo
	Programs       []ProgramInfo

	AssetOutputs      []TxoutputMap
	AssetInputAmount  []AmountMap
	AssetOutputAmount []AmountMap

	Hash string
}

type BlockHead struct {
	Version          uint32
	PrevBlockHash    string
	TransactionsRoot string
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	NextBookKeeper   string
	Program          ProgramInfo

	Hash string
}

type BlockInfo struct {
	Hash         string
	BlockData    *BlockHead
	Transactions []*Transactions
}

type TxInfo struct {
	Hash string
	Hex  string
	Tx   *Transactions
}

type TxoutInfo struct {
	High  uint32
	Low   uint32
	Txout tx.TxOutput
}

type NodeInfo struct {
	State    uint   
	Port     uint16 
	ID       uint64 
	Time     int64
	Version  uint32 
	Services uint64 
	Relay    bool   
	Height   uint64 
	TxnCnt   uint64 
	RxTxnCnt uint64 
}

type ConsensusInfo struct {
	
}

