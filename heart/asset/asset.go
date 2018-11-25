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

