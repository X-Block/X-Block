package common

import (
	. "XBlock/common"
	. "XBlock/common/config"
	"XBlock/core/ledger"
	tx "XBlock/core/transaction"
	. "XBlock/net/httpjsonrpc"
	Err "XBlock/net/httprestful/error"
	. "XBlock/net/protocol"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var node Noder
var pushBlockFlag bool = true

type ApiServer interface {
	Start() error
	Stop()
}

func SetNode(n Noder) {
	node = n
}
func CheckPushBlock() bool {
	return pushBlockFlag
}


func GetConnectionCount(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	if node != nil {
		resp["Result"] = node.GetConnectionCnt()
	}

	return resp
}


