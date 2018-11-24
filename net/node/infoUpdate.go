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

func (node *node) SyncBlk() {
	headerHeight := ledger.DefaultLedger.Store.GetHeaderHeight()
	currentBlkHeight := ledger.DefaultLedger.Blockchain.BlockHeight
	if currentBlkHeight >= headerHeight {
		return
	}
	var dValue int32
	var reqCnt uint32
	var i uint32
	noders := node.local.GetNeighborNoder()
	for _, n := range noders {
		n.RemoveFlightHeightLessThan(currentBlkHeight)
		count := MAXREQBLKONCE - uint32(n.GetFlightHeightCnt())
		dValue = int32(headerHeight - currentBlkHeight - reqCnt)
		for i = 1; i <= count && dValue >= 0; i++ {
			hash := ledger.DefaultLedger.Store.GetHeaderHashByHeight(currentBlkHeight + reqCnt)
			ReqBlkData(n, hash)
			n.StoreFlightHeight(currentBlkHeight + reqCnt)
			reqCnt++
			dValue--
		}
	}
}

func (node *node) SendPingToNbr() {
	noders := node.local.GetNeighborNoder()
	for _, n := range noders {
		t := n.GetLastRXTime()
		if time.Since(t).Seconds() > PERIODUPDATETIME {
			if n.GetState() == ESTABLISH {
				buf, err := NewPingMsg()
				if err != nil {
					log.Error("failed build a new ping message")
				} else {
					go n.Tx(buf)
				}
			}
		}
	}
}

