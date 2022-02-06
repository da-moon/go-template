package value

import "flag"

var _ flag.Value = &String{}

// String provides a flag value that's aware if it has been set.
type String struct {
	v *string
}

// Merge will overlay this value if it has been set.
func (s *String) Merge(onto *string) {
	if s.v != nil {
		*onto = *(s.v)
	}
}

// Set implements the flag.Value interface.
func (s *String) Set(v string) error {
	s.RawSet(v)
	return nil
}

// RawSet sets the underlying value directly
func (s *String) RawSet(v string) {
	if s.v == nil {
		s.v = new(string)
	}
	*(s.v) = v
}

// Get returns the actual underlying value
func (s *String) Get() string {
	var result string
	if s.v != nil {
		result = *(s.v)
	}
	return result
}

// String implements the flag.Value interface.
func (s *String) String() string {
	return s.Get()
}
