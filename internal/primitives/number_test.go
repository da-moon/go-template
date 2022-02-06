package primitives_test

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
	primitives "gitlab.com/tmobile/relic/accelerators/relic-api/internal/primitives"
)

// TestIsInt tests the IsInt function
func TestIsInt(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		data     string
		expected bool
	}{
		{"", false},
		{"nil", false},
		{"1", true},
		{"0", true},
		{"1-", false},
		{"-1", true},
		{"3.14", false},
		{"\u0031", true},
	}
	for _, test := range tests {
		actual := primitives.IsInt(test.data)
		assert.Equal(t, test.expected, actual)
	}
}

// TestIsFloat tests the IsFloat function
func TestIsFloat(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		data     string
		expected bool
	}{
		{"", false},
		{"nil", false},
		{"1", true},
		{"0", true},
		{"1-", false},
		{"-1", true},
		{"-3.15", true},
	}
	for _, test := range tests {
		actual := primitives.IsFloat(test.data)
		assert.Equal(t, test.expected, actual)
	}
}
