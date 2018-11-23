package message

import (
	"XBlock/common"
	"XBlock/common/log"
	"XBlock/core/ledger"
	"XBlock/events"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

