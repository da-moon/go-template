package value

import (
	"flag"
	"strings"
)

var _ flag.Value = (*AppendSlice)(nil)

// AppendSlice implements the flag.Value interface and allows multiple
// calls to the same variable to append a list.
type AppendSlice []string

// Set implements the flag.Value interface.
func (s *AppendSlice) Set(value string) error {
	s.RawSet(value)
	return nil
}

// RawSet sets the underlying value directly
func (s *AppendSlice) RawSet(value string) {
	*s = append(s.Get(), value)
}

// Get returns the actual underlying value
func (s *AppendSlice) Get() []string {
	if *s == nil {
		*s = make([]string, 0, 1)
	}
	return *s
}

// String implements the flag.Value interface.
// [ TODO ] ensure there are no bugs when
// the slice has an empty string
func (s *AppendSlice) String() string {
	return strings.Join(s.Get(), ",")
}
