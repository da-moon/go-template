package primitives_test

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	primitives "github.com/da-moon/go-template/internal/primitives"
	assert "github.com/stretchr/testify/assert"
)

func TestByteSliceFromString(t *testing.T) {
	prevProcs := runtime.GOMAXPROCS(-1)
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(prevProcs)
	src := "bytes! and all these other bytes"
	t.Run("Byte Slice From String Conversion", func(t *testing.T) {
		t.Parallel()
		b := primitives.ByteSliceFromString(src)
		// Should have the same length
		assert.Equal(t, len(b), len(src), "Converted bytes have different length (%d) than the string (%d)", len(b), len(src))
		assert.Equal(t, cap(b), len(src), "Converted bytes have capacity (%d) beyond the length of string (%d)", cap(b), len(src))
		// Should have same content
		assert.Equal(t, string(b), src, "Converted bytes has different value %q than the string %q", string(b), src)
		// Should point to the same data in memory
		sData := (*(*reflect.StringHeader)(unsafe.Pointer(&src))).Data
		bData := (*(*reflect.SliceHeader)(unsafe.Pointer(&b))).Data
		assert.Equal(t, bData, sData, "Converted bytes points to different data %v than the string %v", sData, bData)
	})
	t.Run("ByteSlice From Empty string", func(t *testing.T) {
		t.Parallel()
		got := primitives.ByteSliceFromString("")
		assert.Nil(t, got, "ByteSliceFromString() = %q but want %q", got, "")
	})
}
func TestCutBytes(t *testing.T) {
	expected := []byte("1256")
	actual := primitives.CutBytes([]byte("123456"), 2, 4)
	assert.Equal(t, expected, actual)
}
func TestInsertBytes(t *testing.T) {
	expected := []byte("12abcd3456")
	actual := primitives.InsertBytes([]byte("123456"), 2, []byte("abcd"))
	assert.Equal(t, expected, actual)
}
func TestReplaceBytes(t *testing.T) {
	var tests = []struct {
		expected []byte
		actual   []byte
	}{
		{[]byte("12ab56"), primitives.ReplaceBytes([]byte("123456"), 2, 4, []byte("ab"))},
		{[]byte("12abcd56"), primitives.ReplaceBytes([]byte("123456"), 2, 4, []byte("abcd"))},
		{[]byte("12ab6"), primitives.ReplaceBytes([]byte("123456"), 2, 5, []byte("ab"))},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, test.actual)
	}
}
