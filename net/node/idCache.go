package node

import (
	"XBlock/common"
	"sync"
)

type idCache struct {
	sync.RWMutex
	list map[common.Uint256]bool
}

