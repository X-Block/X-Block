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

type FileStore struct {
	fd   FileData
	file *os.File
	path string
}

func (cs *FileStore) readDB() ([]byte, error) {
	var err error
	cs.file, err = os.OpenFile(cs.path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer cs.closeDB()

	if cs.file != nil {
		data, err := ioutil.ReadAll(cs.file)
		if err != nil {
			return nil, err
		}

		return data, nil

	} else {
		return nil, NewDetailErr(errors.New("[readDB] file handle is nil"), ErrNoCode, "")
	}
}

