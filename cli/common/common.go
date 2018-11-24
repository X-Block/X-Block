package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"XBlock/common/config"

	"github.com/urfave/cli"
)

var (
	Ip   string
	Port string
)

