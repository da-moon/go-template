package testutils

import (
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v6"
)

//nolint:gochecknoglobals
var tmpdir = "/tmp/go-template-test"

//nolint:gochecknoinits
func init() {
	gofakeit.Seed(0)
	err := os.MkdirAll(tmpdir, 0755)
	if err != nil {
		fmt.Printf("Cannot create %s. Reverting to /tmp\n", tmpdir)
		tmpdir = "/tmp"
	}

}
