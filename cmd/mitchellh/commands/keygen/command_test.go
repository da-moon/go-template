package keygen_test

import (
	"strings"
	"testing"

	keygen "github.com/da-moon/go-template/cmd/mitchellh/commands/keygen"
	cli "github.com/mitchellh/cli"
	assert "github.com/stretchr/testify/assert"
)

func TestKeygen(t *testing.T) {
	assert := assert.New(t)
	ui := cli.NewMockUi()
	cmd := keygen.New(ui)
	t.Run("NoTabs", func(t *testing.T) {
		assert.False(strings.ContainsRune(cmd.Help(), '\t'), "help has tabs")
	})
	t.Run("FaultyArgs", func(t *testing.T) {
		assert.NotEqual(cmd.Run([]string{"--foo=123"}), 0)
	})
	t.Run("NoArgs", func(t *testing.T) {
		assert.NotEqual(cmd.Run([]string{}), 0)
	})
}
