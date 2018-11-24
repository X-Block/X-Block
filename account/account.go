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

func NewAccount() (*Account, error) {
	priKey, pubKey, _ := crypto.GenKeyPair()
	signatureRedeemScript, err := contract.CreateSignatureRedeemScript(&pubKey)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "CreateSignatureRedeemScript failed")
	}
	programHash, err := ToCodeHash(signatureRedeemScript)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "ToCodeHash failed")
	}
	return &Account{
		PrivateKey:  priKey,
		PublicKey:   &pubKey,
		ProgramHash: programHash,
	}, nil
}

