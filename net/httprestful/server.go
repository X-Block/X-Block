package httprestful

import (
	. "XBlock/common/config"
	"XBlock/core/ledger"
	"XBlock/events"
	"XBlock/net/httprestful/common"
	. "XBlock/net/httprestful/restful"
	. "XBlock/net/protocol"
	"strconv"
)

const OAUTH_SSUCCESS_CODE = "r0000"

func StartServer(n Noder) {
	common.SetNode(n)
	ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted, SendBlock2NoticeServer)
	func() common.ApiServer {
		rest := InitRestServer(checkAccessToken)
		go rest.Start()
		return rest
	}()
}

func SendBlock2NoticeServer(v interface{}) {

	if len(Parameters.NoticeServerAddr) == 0 || !common.CheckPushBlock() {
		return
	}
	go func() {
		req := make(map[string]interface{})
		req["Height"] = strconv.FormatInt(int64(ledger.DefaultLedger.Blockchain.BlockHeight), 10)
		req = common.GetBlockByHeight(req)

		repMsg, _ := common.PostRequest(req, Parameters.NoticeServerAddr)
		if repMsg[""] == nil {
	
		}
	}()
}

