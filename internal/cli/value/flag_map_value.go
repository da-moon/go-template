package value

import (
	"flag"
	"fmt"
	"strings"

	"github.com/palantir/stacktrace"
)

var _ flag.Value = (*FlagMap)(nil)

// FlagMap is a flag implementation used to provide key=value semantics
// multiple times.
type FlagMap map[string]string

// Set implements the flag.Value interface.
func (h *FlagMap) Set(value string) error {
	idx := strings.Index(value, "=")
	if idx == -1 {
		return stacktrace.NewError("Missing '=' value in argument: %s", value)
	}

	key, value := value[0:idx], value[idx+1:]
	h.RawSet(key, value)
	return nil
}

// RawSet sets the underlying value directly
func (h *FlagMap) RawSet(key, value string) {
	if *h == nil {
		*h = make(map[string]string)
	}
	headers := h.Get()
	headers[key] = value
	*h = headers
}

// Get returns the actual underlying value
func (h *FlagMap) Get() map[string]string {
	if *h == nil {
		*h = make(map[string]string)
	}
	return *h
}

// String implements the flag.Value interface.
// [ TODO ] => maybe pretty print this value
func (h *FlagMap) String() string {
	return fmt.Sprintf("%v", h.Get())
}
