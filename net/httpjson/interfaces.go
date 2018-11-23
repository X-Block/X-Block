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

