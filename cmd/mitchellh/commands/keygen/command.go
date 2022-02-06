package keygen

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	primitives "github.com/da-moon/go-template/internal/primitives"
	cli "github.com/mitchellh/cli"
)

const (
	entrypoint = "keygen"
	length     = 32
)

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.init()
	return c
}

type cmd struct {
	UI        cli.Ui
	flags     *Flags
	help      string
	synopsis  string
	testStdin io.Reader
}

func (c *cmd) init() {
	c.UI = &cli.PrefixedUi{
		OutputPrefix: "",
		InfoPrefix:   "",
		ErrorPrefix:  "",
		Ui:           c.UI,
	}
	c.flags = &Flags{}
	c.flags.init()

	c.synopsis = synopsis
	c.help = c.flags.Usage()
	c.flags.Value().Usage = func() { c.UI.Info(c.Help()) }
}

// Run ...
func (c *cmd) Run(args []string) int {
	if c.flags == nil {
		c.UI.Error("underlying flag struct was nil")
		return 1
	}
	flags := c.flags.Value()
	err := flags.Parse(args)
	if err != nil {
		return 1
	}
	args = flags.Args()
	// TODO fix this
	if len(args) != 0 {
		c.UI.Error("this subcommand takes no arguments")
		return 1
	}
	key := make([]byte, length)
	n, err := rand.Reader.Read(key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not read random data: %s", err))
		return 1
	}
	if n != length {
		c.UI.Error("could not read enough entropy. Generate more entropy!")
		return 1
	}
	hexResult := hex.EncodeToString(key)
	baseResult := base64.StdEncoding.EncodeToString(key)
	hexEncode := c.flags.Hex()
	baseEncode := c.flags.Base64()
	if hexEncode && baseEncode {
		type response struct {
			Hex    string `json:"hex"`
			Base64 string `json:"base64"`
		}
		result, err := primitives.IndentedJSON(response{
			Hex:    hexResult,
			Base64: baseResult,
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("err:%v", err))
			return 1
		}
		c.UI.Output(result.String())
		return 0
	} else if hexEncode {
		c.UI.Output(hexResult)
		return 0
	} else if baseEncode {
		c.UI.Output(baseResult)
		return 0
	}
	c.UI.Error("you must choose encoding scheme!")
	return 1
}

// Synopsis ...
func (c *cmd) Synopsis() string {
	return strings.TrimSpace(c.synopsis)
}

// Help ...
func (c *cmd) Help() string {
	return strings.TrimSpace(c.help)
}

const synopsis = "Generates a new encryption key"

const help = `
Usage: go-template keygen
  Generates a new 32 bytes long random encryption key that can be used to for encrypting data.
  in case both base64 and hex output is requested, the result will be returned as a json
  reply, containing both values.
`
