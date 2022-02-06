package primitives_test

import (
	"strings"
	"testing"

	primitives "github.com/da-moon/go-template/internal/primitives"
	assert "github.com/stretchr/testify/assert"
)

// TestMapString tests the MapString function
func TestMapString(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []string
		expected []string
	}{
		{[]string{"foo", "bar"}, []string{"FOO", "BAR"}},
		{[]string{"foo", "\u0062\u0061\u0072"}, []string{"FOO", "BAR"}},
		{[]string{}, []string{}},
	}
	for _, test := range tests {
		actual := primitives.MapString(test.s, strings.ToUpper)
		assert.True(t, primitives.EqSlices(&actual, &test.expected), "Expected MapString(%q, fn) to be %q, got %v", test.s, test.expected, actual)
	}
}

// TestMapInt tests the MapInt function
func TestMapInt(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []int
		expected []int
	}{
		{[]int{0, 1, 2}, []int{0, 2, 4}},
		{[]int{-1}, []int{-2}},
		{[]int{}, []int{}},
	}
	for _, test := range tests {
		actual := primitives.MapInt(test.s, func(i int) int {
			return i * 2
		})
		assert.True(t, primitives.EqSlices(&actual, &test.expected), "Expected MapInt(%q, fn) to be %q, got %v", test.s, test.expected, actual)
	}
}

// TestFilterString tests the FilterString function
func TestFilterString(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []string
		expected []string
	}{
		{[]string{"foo", "bar", "baz"}, []string{"bar", "baz"}},
		{[]string{"foo", "\u0062\u0061\u0072", "baz"}, []string{"bar", "baz"}},
		{[]string{"a", "ab", "abc"}, []string{}},
		{[]string{}, []string{}},
	}
	for _, test := range tests {
		actual := primitives.FilterString(test.s, func(s string) bool {
			return strings.HasPrefix(s, "ba")
		})
		assert.True(t, primitives.EqSlices(&actual, &test.expected), "Expected FilterString(%q, fn) to be %q, got %v", test.s, test.expected, actual)
	}
}

// TestFilterInt tests the FilterInt function
func TestFilterInt(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []int
		expected []int
	}{
		{[]int{0, 2, 4}, []int{0, 2, 4}},
		{[]int{}, []int{}},
		{[]int{2, 4, 1}, []int{2, 4}},
		{[]int{1}, []int{}},
		{[]int{-2, 4}, []int{-2, 4}},
	}
	for _, test := range tests {
		actual := primitives.FilterInt(test.s, func(i int) bool {
			return i%2 == 0
		})
		assert.True(t, primitives.EqSlices(&actual, &test.expected), "Expected FilterInt(%q, fn) to be %q, got %v", test.s, test.expected, actual)
	}
}

// TestAllString tests the AllString function
func TestAllString(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []string
		expected bool
	}{
		{[]string{"boo", "bar", "baz"}, true},
		{[]string{"boo", "\u0062\u0061\u0072", "baz"}, true},
		{[]string{"foo", "bar", "baz"}, false},
		{[]string{}, true},
	}
	for _, test := range tests {
		actual := primitives.AllString(test.s, func(s string) bool {
			return strings.HasPrefix(s, "b")
		})
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestAllInt tests the AllInt function
func TestAllInt(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []int
		expected bool
	}{
		{[]int{0, 2, 4}, true},
		{[]int{}, true},
		{[]int{2, 4, 1}, false},
		{[]int{1}, false},
		{[]int{-2, 4}, true},
	}
	for _, test := range tests {
		actual := primitives.AllInt(test.s, func(i int) bool {
			return i%2 == 0
		})
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestAnyString tests the AnyString function
func TestAnyString(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []string
		expected bool
	}{
		{[]string{"foo", "\u0062\u0061\u0072", "baz"}, true},
		{[]string{"boo", "bar", "baz"}, false},
		{[]string{"foo", "far", "baz"}, true},
	}
	for _, test := range tests {
		actual := primitives.AnyString(test.s, func(s string) bool {
			return strings.HasPrefix(s, "f")
		})
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestAnyInt tests the AnyInt function
func TestAnyInt(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s        []int
		expected bool
	}{
		{[]int{0, 2, 4}, true},
		{[]int{-2, 4}, true},
		{[]int{1}, false},
		{[]int{}, false},
	}
	for _, test := range tests {
		actual := primitives.AnyInt(test.s, func(i int) bool {
			return i%2 == 0
		})
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestIndexi tests the Indexi function
func TestIndexi(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		haystack []string
		needle   string
		expected int
	}{
		{[]string{"FOO", "bar"}, "bar", 1},
		{[]string{"FoO", "bar"}, "foo", 0},
		{[]string{"foo", "bar"}, "blah", -1},
	}
	for _, test := range tests {
		actual := primitives.Indexi(test.haystack, test.needle)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestStringIndex tests the Index function with strings slice
func TestStringIndex(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		haystack []string
		needle   string
		expected int
	}{
		{[]string{"foo", "bar"}, "foo", 0},
		{[]string{"foo", "bar"}, "bar", 1},
		{[]string{"foo", "bar"}, "\u0062\u0061\u0072", 1},
		{[]string{"foo", "bar"}, "", -1},
		{[]string{"foo", "bar"}, "blah", -1},
	}
	for _, test := range tests {
		actual := primitives.Index(&test.haystack, test.needle)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestIntIndex tests the Index function with ints slice
func TestIntIndex(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		haystack []int
		needle   int
		expected int
	}{
		{[]int{1, 2}, 1, 0},
		{[]int{1, 2}, 2, 1},
		{[]int{1, 2}, 0, -1},
	}
	for _, test := range tests {
		actual := primitives.Index(&test.haystack, test.needle)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestFloatIndex tests the Index function with ints slice
func TestFloatIndex(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		haystack []float64
		needle   float64
		expected int
	}{
		{[]float64{1, 3.14}, 1, 0},
		{[]float64{1, 3.14}, 3.14, 1},
		{[]float64{1, 3}, 0, -1},
		{[]float64{1, 2}, 2, 1},
	}
	for _, test := range tests {
		actual := primitives.Index(&test.haystack, test.needle)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestEqStringSlices tests the EqSlices function with strings slice
func TestEqStringSlice(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		a        []string
		b        []string
		expected bool
	}{
		{[]string{"foo", "bar"}, []string{"foo", "bar"}, true},
		{[]string{"foo", "bar"}, []string{"bar", "foo"}, false},
		{[]string{"foo", "bar"}, []string{"bar"}, false},
		{[]string{"foo", "bar"}, []string{"\x66\x6f\x6f", "bar"}, true},
	}
	for _, test := range tests {
		actual := primitives.EqSlices(&test.a, &test.b)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestEqIntSlices tests the EqSlices function with strings slice
func TestEqIntSlice(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		a        []int
		b        []int
		expected bool
	}{
		{[]int{1, 2}, []int{1, 2}, true},
		{[]int{1, 2}, []int{2, 1}, false},
		{[]int{1, 2}, []int{1}, false},
		{[]int{1, 2}, []int{1, 2, 1}, false},
	}
	for _, test := range tests {
		actual := primitives.EqSlices(&test.a, &test.b)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestInStringSlice tests the InSlice function with strings slice
func TestInStringSlice(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		haystack []string
		needle   string
		expected bool
	}{
		{[]string{"foo", "bar"}, "foo", true},
		{[]string{"foo", "bar"}, "", false},
		{[]string{"foo", "bar"}, "f", false},
	}
	for _, test := range tests {
		actual := primitives.InSlice(test.needle, &test.haystack)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}

// TestInIntSlice tests the InSlice function with ints slice
func TestInIntSlice(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		haystack []int
		needle   int
		expected bool
	}{
		{[]int{0, 1, 2}, 2, true},
		{[]int{0, 1, 2}, 3, false},
	}
	for _, test := range tests {
		actual := primitives.InSlice(test.needle, &test.haystack)
		assert.Equal(t, test.expected, actual, "expected value '%v' | actual : '%v'", test.expected, actual)
	}
}
