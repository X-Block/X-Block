package httpjson

import (
	"XBlock/common/log"
	. "XBlock/common/config"
	"net/http"
	"strconv"
)


func StartLocalServer() {
	log.Debug()
	http.HandleFunc(LocalDir, Handle)

	HandleFunc("getneighbor", getNeighbor)
	HandleFunc("getnodestate", getNodeState)
	HandleFunc("startconsensus", startConsensus)
	HandleFunc("stopconsensus", stopConsensus)
	HandleFunc("sendsampletransaction", sendSampleTransaction)
	HandleFunc("setdebuginfo", setDebugInfo)

	
