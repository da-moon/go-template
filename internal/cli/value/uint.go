package value

import (
	"flag"
	"fmt"
	"strconv"
)

var _ flag.Value = &Uint{}

// Uint provides a flag value that's aware if it has been set.
type Uint struct {
	v *uint
}

// Merge will overlay this value if it has been set.
func (u *Uint) Merge(onto *uint) {
	if u.v != nil {
		*onto = *(u.v)
	}
}

// Set implements the flag.Value interface.
func (u *Uint) Set(v string) error {
	parsed, err := strconv.ParseUint(v, 0, 64)
	u.RawSet((uint)(parsed))
	return err
}

// RawSet sets the underlying value directly
func (u *Uint) RawSet(v uint) {
	if u.v == nil {
		u.v = new(uint)
	}
	*(u.v) = v
}

// Get returns the actual underlying value
func (u *Uint) Get() uint {
	var result uint
	if u.v != nil {
		result = *(u.v)
	}
	return result
}

// String implements the flag.Value interface.
func (u *Uint) String() string {
	return fmt.Sprintf("%v", u.Get())
}
