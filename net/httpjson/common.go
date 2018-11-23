package httpjson

import (
	. "XBlock/common"
	"XBlock/common/log"
	"XBlock/consensus/dbft"
	"XBlock/core/ledger"
	. "XBlock/core/transaction"
	tx "XBlock/core/transaction"
	"XBlock/core/validation"
	. "XBlock/net/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

