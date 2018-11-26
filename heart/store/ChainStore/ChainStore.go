package ChainStore

import (
	. "XBlock/common"
	"XBlock/common/log"
	"XBlock/common/serialization"
	. "XBlock/core/asset"
	"XBlock/core/contract/program"
	. "XBlock/core/ledger"
	. "XBlock/core/store"
	. "XBlock/core/store/LevelDBStore"
	tx "XBlock/core/transaction"
	"XBlock/core/transaction/payload"
	"XBlock/core/validation"
	"XBlock/crypto"
	. "XBlock/errors"
	"XBlock/events"
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"XBlock/core/account"
)


var (
	ErrDBNotFound = errors.New("leveldb: not found")
)

type ChainStore struct {
	st IStore

	headerIndex map[uint32]Uint256
	blockCache  map[Uint256]*Block
	headerCache map[Uint256]*Header

	currentBlockHeight uint32
	storedHeaderCount  uint32

	mu sync.RWMutex

	disposed bool
}

func init() {
}

func NewStore() IStore {
	ldbs, _ := NewLevelDBStore("Chain")

	return ldbs
}

func NewLedgerStore() ILedgerStore {
	cs, _ := NewChainStore("Chain")

	return cs
}

