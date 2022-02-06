package main

import (
	"os"

	keygen "github.com/da-moon/go-template/cmd/mitchellh/commands/keygen"
	version "github.com/da-moon/go-template/cmd/mitchellh/commands/version"
	cli "github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	Commands = map[string]cli.CommandFactory{
		"keygen": func() (cli.Command, error) {
			return keygen.New(ui), nil
		},
		"version": func() (cli.Command, error) {
			return version.New(ui), nil
		},
	}
}
