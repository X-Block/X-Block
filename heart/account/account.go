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

func NewAccountState(programHash common.Uint160, balances map[common.Uint256]common.Fixed64) *AccountState {
	var accountState AccountState
	accountState.ProgramHash = programHash
	accountState.Balances = balances
	accountState.IsFrozen = false
	return &accountState
}

