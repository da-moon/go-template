package value

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var _ flag.Value = &Bool{}

// Bool provides a flag value that's aware if it has been set.
type Bool struct {
	v *bool
}

// IsBoolFlag is an optional method of the flag.Value
// interface which marks this value as boolean when
// the return value is true. See flag.Value for details.
func (b *Bool) IsBoolFlag() bool {
	return true
}

// Merge will overlay this value if it has been set.
func (b *Bool) Merge(onto *bool) {
	if b.v != nil {
		*onto = *(b.v)
	}
}

// Set implements the flag.Value interface.
func (b *Bool) Set(v string) error {
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	b.RawSet(parsed)
	return nil
}

// RawSet sets the underlying value directly
func (b *Bool) RawSet(v bool) {
	if b.v == nil {
		b.v = new(bool)
	}
	*(b.v) = v
}

// Get returns the actual underlying value
func (b *Bool) Get() bool {
	var result bool
	if b.v != nil {
		result = *(b.v)
	}
	return result
}

// String implements the flag.Value interface.
func (b *Bool) String() string {
	result := b.Get()
	return strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", result)))
}
