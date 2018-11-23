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

func RegistRpcNode(n Noder) {
	if node == nil {
		node = n
	}
}

func RegistDbftService(d *dbft.DbftService) {
	if dBFT == nil {
		dBFT = d
	}
}


func HandleFunc(pattern string, handler func([]interface{}) map[string]interface{}) {
	mainMux.Lock()
	defer mainMux.Unlock()
	mainMux.m[pattern] = handler
}


func SetDefaultFunc(def func(http.ResponseWriter, *http.Request)) {
	mainMux.defaultFunction = def
}

func Handle(w http.ResponseWriter, r *http.Request) {
	mainMux.RLock()
	defer mainMux.RUnlock()
	if r.Method != "POST" {
		if mainMux.defaultFunction != nil {
			log.Info("HTTP JSON RPC Handle - Method!=\"POST\"")
			mainMux.defaultFunction(w, r)
			return
		} else {
			log.Warn("HTTP JSON RPC Handle - Method!=\"POST\"")
			return
		}
	}

	if r.Body == nil {
		if mainMux.defaultFunction != nil {
			log.Info("HTTP JSON RPC Handle - Request body is nil")
			mainMux.defaultFunction(w, r)
			return
		} else {
			log.Warn("HTTP JSON RPC Handle - Request body is nil")
			return
		}
	}

	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("HTTP JSON RPC Handle - ioutil.ReadAll: ", err)
		return
	}
	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Error("HTTP JSON RPC Handle - json.Unmarshal: ", err)
		return
	}

	
	function, ok := mainMux.m[request["method"].(string)]
	if ok {
		response := function(request["params"].([]interface{}))
		data, err := json.Marshal(map[string]interface{}{
			"jsonpc": "2.0",
			"result": response["result"],
			"id":     request["id"],
		})
		if err != nil {
			log.Error("HTTP JSON RPC Handle - json.Marshal: ", err)
			return
		}
		w.Write(data)
	} else {
		
		log.Warn("HTTP JSON RPC Handle - No function to call for ", request["method"])
		data, err := json.Marshal(map[string]interface{}{
			"result": nil,
			"error": map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
				"data":    "The called method was not found on the server",
			},
			"id": request["id"],
		})
		if err != nil {
			log.Error("HTTP JSON RPC Handle - json.Marshal: ", err)
			return
		}
		w.Write(data)
	}
}

func responsePacking(result interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": result,
	}
	return resp
}


func Call(address string, method string, id interface{}, params []interface{}) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Marshal JSON request: %v\n", err)
		return nil, err
	}
	resp, err := http.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "POST request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GET response: %v\n", err)
		return nil, err
	}

	return body, nil
}

