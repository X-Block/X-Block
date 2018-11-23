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

