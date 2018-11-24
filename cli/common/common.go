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

func NewIpFlag() cli.Flag {
	return cli.StringFlag{
		Name:        "ip",
		Usage:       "node's ip address",
		Value:       "localhost",
		Destination: &Ip,
	}
}

