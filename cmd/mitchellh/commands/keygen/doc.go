// Package keygen creates a subcommand that can help with
// generating random strings.
package keygen

import cli "github.com/mitchellh/cli"

var _ cli.Command = &cmd{}
