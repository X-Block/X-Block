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

func openWallet(name string, passwd []byte) account.Client {
	if name == account.WalletFileName {
		fmt.Println("Using default wallet: ", account.WalletFileName)
	}
	wallet := account.Open(name, passwd)
	if wallet == nil {
		fmt.Println("Failed to open wallet: ", name)
		os.Exit(1)
	}
	return wallet
}

func getUintHash(programHashStr, assetHashStr string) (Uint160, Uint256, error) {
	programHashHex, err := hex.DecodeString(programHashStr)
	if err != nil {
		fmt.Println("Decoding program hash string failed")
		return Uint160{}, Uint256{}, err
	}
	var programHash Uint160
	if err := programHash.Deserialize(bytes.NewReader(programHashHex)); err != nil {
		fmt.Println("Deserialization program hash failed")
		return Uint160{}, Uint256{}, err
	}
	assetHashHex, err := hex.DecodeString(assetHashStr)
	if err != nil {
		fmt.Println("Decoding asset hash string failed")
		return Uint160{}, Uint256{}, err
	}
	var assetHash Uint256
	if err := assetHash.Deserialize(bytes.NewReader(assetHashHex)); err != nil {
		fmt.Println("Deserialization asset hash failed")
		return Uint160{}, Uint256{}, err
	}
	return programHash, assetHash, nil
}

