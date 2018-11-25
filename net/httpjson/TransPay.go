package httpjsonrpc

import (
	. "XBlock/common"
	"XBlock/core/asset"
	. "XBlock/core/contract"
	. "XBlock/core/transaction"
	"XBlock/core/transaction/payload"
	"bytes"
)

type Payload interface {
	Data() []byte
}


type BookKeeping struct {
}

func (dc *BookKeeping) Data() []byte {
	return []byte{0}
}


type FunctionCode struct {
	Code           string
	ParameterTypes string
	ReturnTypes    string
}

type DeployCode struct {
	Code        *FunctionCode
	Name        string
	CodeVersion string
	Author      string
	Email       string
	Description string
}

func (dc *DeployCode) Data() []byte {
	return []byte{0}
}


type IssueAsset struct {
}

func (a *IssueAsset) Data() []byte {
	return []byte{0}
}

type Issuer struct {
	X, Y string
}


type RegisterAsset struct {
	Asset      *asset.Asset
	Amount     Fixed64
	Issuer     Issuer
	Controller string
}

func (a *RegisterAsset) Data() []byte {
	return []byte{0}
}


