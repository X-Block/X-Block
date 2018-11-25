package asset

import (
	"XBlock/common/serialization"
	. "XBlock/errors"
	"errors"
	"io"
)


type AssetType byte

const (
	Currency AssetType = 0x00
	Share    AssetType = 0x01
	Invoice  AssetType = 0x10
	Token    AssetType = 0x11
)

type AssetRecordType byte

type Asset struct {
	Name       string
	Precision  byte
	AssetType  AssetType
	RecordType AssetRecordType
}


func (a *Asset) Serialize(w io.Writer) error {
	err := serialization.WriteVarString(w, a.Name)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], Name serialize failed.")
	}
	_, err = w.Write([]byte{byte(a.Precision)})
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], Precision serialize failed.")
	}
	_, err = w.Write([]byte{byte(a.AssetType)})
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], AssetType serialize failed.")
	}
	_, err = w.Write([]byte{byte(a.RecordType)})
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], RecordType serialize failed.")
	}
	return nil
}


