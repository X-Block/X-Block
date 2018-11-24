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

