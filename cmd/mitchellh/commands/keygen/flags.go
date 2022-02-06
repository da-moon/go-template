package keygen

import (
	"bytes"
	"fmt"

	flagset "github.com/da-moon/go-template/internal/cli/flagset"
	value "github.com/da-moon/go-template/internal/cli/value"
)

type Flags struct {
	*flagset.FlagSet
	base64 value.Bool
	hex    value.Bool
}

func (f *Flags) init() {
	f.FlagSet = flagset.New(entrypoint, help)
	f.FlagSet.Var(&f.base64, "base64",
		"encodes the result as a base64 string")
	f.FlagSet.Var(&f.hex, "hex",
		"encodes the result as a hex string ")
}

func (f *Flags) Base64() bool {
	return f.base64.Get()
}
func (f *Flags) Hex() bool {
	return f.hex.Get()
}
func (f *Flags) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "\nbase64\t\t:\t%s", f.base64.String())
	fmt.Fprintf(&buf, "\nhex\t:\t%s", f.hex.String())
	return buf.String()
}
