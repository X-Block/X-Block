package account

import (
	"XBlock/common"
	"io"
	"bytes"
	"XBlock/common/serialization"
)

type AccountState struct {
	ProgramHash common.Uint160
	IsFrozen bool
	Balances map[common.Uint256]common.Fixed64
}

