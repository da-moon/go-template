package version

import (
	"strings"

	buildver "github.com/da-moon/go-template/build/version"
	cli "github.com/mitchellh/cli"
)

func New(ui cli.Ui) *cmd {
	c := &cmd{
		UI: ui,
	}
	return c
}

type cmd struct {
	UI       cli.Ui
	help     string
	synopsis string
}

// Run ...
func (c *cmd) Run(_ []string) int {
	build := buildver.New()
	c.UI.Output(build.ToString())
	return 0
}
func (c *cmd) Synopsis() string {
	return strings.TrimSpace(c.synopsis)
}

// Help ...
func (c *cmd) Help() string {
	return strings.TrimSpace(c.help)
}

const synopsis = "go-template immutable build information ."

const help = `
Usage: go-template version

returns current release buildver.
`
