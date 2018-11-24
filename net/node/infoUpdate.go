package node

import (
	"XBlock/common/config"
	"XBlock/common/log"
	"XBlock/core/ledger"
	. "XBlock/net/message"
	. "XBlock/net/protocol"
	"fmt"
	"net"
	"strconv"
	"time"
)

func keepAlive(from *Noder, dst *Noder) {

}

func (node *node) GetBlkHdrs() {
	if node.local.GetNbrNodeCnt() < MINCONNCNT {
		return
	}

	noders := node.local.GetNeighborNoder()
	for _, n := range noders {
		if uint64(ledger.DefaultLedger.Store.GetHeaderHeight()) < n.GetHeight() {
			if n.LocalNode().IsSyncFailed() == false {
				SendMsgSyncHeaders(n)
				n.StartRetryTimer()
				break
			}
		}
	}
}

