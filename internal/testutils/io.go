package testutils

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	primitives "github.com/da-moon/go-template/internal/primitives"
)

// RWTestCase define test cases used in testing Write([]byte) and Read([]byte)
// for unit , use values in primitives package
type RWTestCase struct {
	Size                     int
	Unit                     int
	Expect                   []byte
	CustomWriteTester        func(t *testing.T, w Writer)
	CustomWriteToTester      func(t *testing.T, w WriterTo)
	CustomReaderWriterTester func(t *testing.T, r ReaderWriter)
}

// Writer ...
type Writer interface {
	Write(p []byte) (int, error)
	Bytes() []byte
	Len() int
	Reset()
}

// ReaderWriter ...
type ReaderWriter interface {
	Writer
	Read(p []byte) (n int, err error)
}

// WriterTo ...
type WriterTo interface {
	Writer
	WriteTo(w io.Writer) (int64, error)
}

// Reader ...
type Reader interface {
	Read(p []byte) (n int, err error)
	Len() int
}

// GenerateRWTests generates test for
// to help with testing io.writer and readers
func GenerateRWTests() []RWTestCase {
	writeSizes := []int{1, 2, 4, 8, 16}
	// writeSizes := []int{1, 2, 4, 8, 16, 32}
	// writeSizes := []int{1}
	// unit must be equal or larger than 64
	units := []int{64, primitives.Ki, primitives.Mi}
	// units := []int{64}
	tests := make([]RWTestCase, 0)
	// setting test values
	for _, unit := range units {
		for _, size := range writeSizes {
			randomizer := NewRandomReader(size * unit)
			expect := make([]byte, size*unit)
			randomizer.Read(expect) //nolint:errcheck
			test := RWTestCase{
				Size:   size,
				Unit:   unit,
				Expect: expect,
			}
			tests = append(tests, test)
		}
	}
	return tests
}

// ─── TEMPORARY FILES AND FOLDERS ────────────────────────────────────────────────
//nolint:gochecknoglobals
var noCleanup = strings.ToLower(os.Getenv("TEST_NOCLEANUP")) == "true"

func TempDir(t testing.TB, name string) string {
	if t == nil {
		panic("argument t must be non-nil")
	}
	name = t.Name() + "-" + name
	name = strings.ReplaceAll(name, "/", "_")
	d, err := ioutil.TempDir(tmpdir, name)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	t.Cleanup(func() {
		if noCleanup {
			t.Logf("skipping cleanup because TEST_NOCLEANUP was enabled")
			return
		}
		os.RemoveAll(d)
	})
	return d
}
func TempFile(t testing.TB, name string) *os.File {
	if t == nil {
		panic("argument t must be non-nil")
	}
	name = t.Name() + "-" + name
	name = strings.ReplaceAll(name, "/", "_")
	f, err := ioutil.TempFile(tmpdir, name)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	t.Cleanup(func() {
		if noCleanup {
			t.Logf("skipping cleanup because TEST_NOCLEANUP was enabled")
			return
		}
		os.Remove(f.Name())
	})
	return f
}
