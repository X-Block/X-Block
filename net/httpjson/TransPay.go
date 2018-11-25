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


type TransferAsset struct {
}

func (a *TransferAsset) Data() []byte {
	return []byte{0}
}

type Record struct {
	RecordType string
	RecordData string
}

func (a *Record) Data() []byte {
	return []byte{0}
}

type PrivacyPayload struct {
	PayloadType uint8
	Payload     string
	EncryptType uint8
	EncryptAttr string
}

func (a *PrivacyPayload) Data() []byte {
	return []byte{0}
}

func TransPay(p Payload) Payload {
	switch object := p.(type) {
	case *payload.BookKeeping:
	case *payload.BookKeeper:
	case *payload.IssueAsset:
	case *payload.TransferAsset:
	case *payload.DeployCode:
		obj := new(DeployCode)
		obj.Code.Code = ToHexString(object.Code.Code)
		obj.Code.ParameterTypes = ToHexString(ContractParameterTypeToByte(object.Code.ParameterTypes))
		obj.Code.ReturnTypes = ToHexString(ContractParameterTypeToByte(object.Code.ReturnTypes))
		obj.Name = object.Name
		obj.CodeVersion = object.CodeVersion
		obj.Author = object.Author
		obj.Email = object.Email
		obj.Description = object.Description
		return obj
	

	}
	return nil
}
