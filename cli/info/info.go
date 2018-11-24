package info

import (
	"fmt"
	"os"

	. "XBlock/cli/common"
	"XBlock/net/httpjson"

	"github.com/urfave/cli"
)

func infoAction(c *cli.Context) (err error) {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	blockhash := c.String("blockhash")
	txhash := c.String("txhash")
	bestblockhash := c.Bool("bestblockhash")
	height := c.Int("height")
	blockcount := c.Bool("blockcount")
	connections := c.Bool("connections")
	neighbor := c.Bool("neighbor")
	state := c.Bool("state")
	version := c.Bool("nodeversion")

	var resp []byte
	var output [][]byte
	if height != -1 {
		resp, err = httpjsonrpc.Call(Address(), "getblock", 0, []interface{}{height})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if c.String("blockhash") != "" {
		resp, err = httpjsonrpc.Call(Address(), "getblock", 0, []interface{}{blockhash})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if bestblockhash {
		resp, err = httpjsonrpc.Call(Address(), "getbestblockhash", 0, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if blockcount {
		resp, err = httpjsonrpc.Call(Address(), "getblockcount", 0, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if connections {
		resp, err = httpjsonrpc.Call(Address(), "getconnectioncount", 0, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if neighbor {
		resp, err := httpjsonrpc.Call(Address(), "getneighbor", 0, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if state {
		resp, err := httpjsonrpc.Call(Address(), "getnodestate", 0, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if txhash != "" {
		resp, err = httpjsonrpc.Call(Address(), "getrawtransaction", 0, []interface{}{txhash})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)
	}

	if version {
		resp, err = httpjsonrpc.Call(Address(), "getversion", 0, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		output = append(output, resp)

	}
	for _, v := range output {
		FormatOutput(v)
	}

	return nil
}

