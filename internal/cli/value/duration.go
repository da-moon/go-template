package value

import (
	"flag"
	"time"
)

var _ flag.Value = &Duration{}

// Duration provides a flag value that's aware if it has been set.
type Duration struct {
	v *time.Duration
}

// Merge will overlay this value if it has been set.
func (d *Duration) Merge(onto *time.Duration) {
	if d.v != nil {
		*onto = *(d.v)
	}
}

// Set implements the flag.Value interface.
func (d *Duration) Set(v string) error {
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return err
	}
	d.RawSet(parsed)
	return nil
}

// RawSet sets the underlying value directly
func (d *Duration) RawSet(v time.Duration) {
	if d.v == nil {
		d.v = new(time.Duration)
	}
	*(d.v) = v
}

// Get returns the actual underlying value
func (d *Duration) Get() time.Duration {
	var result time.Duration
	if d.v != nil {
		result = *(d.v)
	}
	return result
}

// String implements the flag.Value interface.
func (d *Duration) String() string {
	return d.Get().String()
}
