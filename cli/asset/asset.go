package asset

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	. "XBlock/cli/common"
	"XBlock/account"
	. "XBlock/common"
	. "XBlock/core/asset"
	"XBlock/core/contract"
	"XBlock/core/signature"
	"XBlock/core/transaction"
	"XBlock/net/httpjsonrpc"

	"github.com/urfave/cli"
)

const (
	RANDBYTELEN    = 4
	REFERTXHASHLEN = 64
)

func newContractContextWithoutProgramHashes(data signature.SignableData) *contract.ContractContext {
	return &contract.ContractContext{
		Data:       data,
		Codes:      make([][]byte, 1),
		Parameters: make([][][]byte, 1),
	}
}

