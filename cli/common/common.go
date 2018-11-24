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

func NewPortFlag() cli.Flag {
	return cli.StringFlag{
		Name:        "port",
		Usage:       "node's RPC port",
		Value:       strconv.Itoa(config.Parameters.HttpLocalPort),
		Destination: &Port,
	}
}

func Address() string {
	address := "http://" + Ip + ":" + Port
	return address
}

