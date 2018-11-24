package account

import (
	ct "XBlock/core/contract"
	. "XBlock/errors"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type FileData struct {
	PublicKeyHash       string
	PrivateKeyEncrypted string
	Address             string
	ScriptHash          string
	RawData             string
	PasswordHash        string
	IV                  string
	MasterKey           string
}

