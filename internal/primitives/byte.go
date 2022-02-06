package primitives

import (
	"fmt"

	"github.com/modern-go/reflect2"
	"github.com/palantir/stacktrace"
)

// Binary prefixes for common use.
const (
	_ = iota
	// Ki ...
	Ki = 1 << (10 * iota)
	// Mi ...
	Mi
	// Gi ...
	Gi
	// Ti ...
	Ti
	// Pi ...
	Pi
	// Ei ...
	Ei
)

// errors
var (
	// ErrByteSliceTooLarge ...
	ErrByteSliceTooLarge = stacktrace.NewError("ByteSlice: too large")
)

// MakeByteSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func MakeByteSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrByteSliceTooLarge)
		}
	}()
	return make([]byte, n)
}

// ByteCountDecimal returns a string represenation of bytes with a base of 10
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountBinary returns a string representaion of bytes with a base of 2
func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// ByteSliceFromString converts a string into a slice of bytes without performing a copy.
// [NOTE] => This is an unsafe operation and may lead to problems if the bytes are changed.
func ByteSliceFromString(str string) []byte {
	return reflect2.UnsafeCastString(str)
}

// CutBytes elements from slice for a given range
func CutBytes(a []byte, from, to int) []byte {
	copy(a[from:], a[to:])
	a = a[:len(a)-to+from]
	return a
}

// InsertBytes new slice at specified position
func InsertBytes(a []byte, i int, b []byte) []byte {
	a = append(a, make([]byte, len(b))...)
	copy(a[i+len(b):], a[i:])
	copy(a[i:i+len(b)], b)
	return a
}

// ReplaceBytes function unlike bytes.Replace allows you to specify range
func ReplaceBytes(a []byte, from, to int, new []byte) []byte {
	lenDiff := len(new) - (to - from)
	if lenDiff > 0 {
		// Extend if new segment bigger
		a = append(a, make([]byte, lenDiff)...)
		copy(a[to+lenDiff:], a[to:])
		copy(a[from:from+len(new)], new)
		return a
	}
	if lenDiff < 0 {
		copy(a[from:], new)
		copy(a[from+len(new):], a[to:])
		return a[:len(a)+lenDiff]
	}
	// same size
	copy(a[from:], new)
	return a
}
