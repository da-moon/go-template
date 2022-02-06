// Package version returns Reproducible and immutable build
// version information
// https://blog.macrium.com/reproducible-and-immutable-builds-can-improve-trust-in-the-software-supply-chain-after-the-17b9117c6e7b
package version

import cli "github.com/mitchellh/cli"

var _ cli.Command = &cmd{}
