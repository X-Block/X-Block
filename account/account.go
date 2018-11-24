package account

import (
	. "XBlock/common"
	"XBlock/core/contract"
	"XBlock/crypto"
	. "XBlock/errors"
	"errors"
)

type Account struct {
	PrivateKey  []byte
	PublicKey   *crypto.PubKey
	ProgramHash Uint160
}

