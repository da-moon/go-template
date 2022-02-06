package testutils

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//nolint:gochecknoglobals
var (
	ignore = []string{
		"runtime.goexit",
		"runtime.main",
	}
	pool = sync.Pool{
		New: func() interface{} {
			result := make([]uintptr, 64)
			return result
		},
	}
)

// CapturedStacktrace ...
type CapturedStacktrace string

// Stacktrace - captures a stacktrace of the current goroutine
// github.com/hashicorp/go-hclog
func Stacktrace() CapturedStacktrace {
	pcs := pool.Get().([]uintptr)
	defer pool.Put(pcs) //nolint
	buffer := new(bytes.Buffer)
	for {
		// Skip the call to runtime.Counters so that the
		// program counters start fresh at the method that invoked
		// CapturedStacktrace().
		n := runtime.Callers(1, pcs)
		if n < cap(pcs) {
			pcs = pcs[:n]
			break
		}
		// if the counter slice is too-short , do not put it back into the pool
		// in case of consistently taking deep stacktraces, this would let
		// the pool adjust.
		pcs = make([]uintptr, len(pcs)*2)
	}
	i := 0
	frames := runtime.CallersFrames(pcs)
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		for _, prefix := range ignore {
			if strings.HasPrefix(frame.Function, prefix) {
				continue
			}
		}
		if i != 0 {
			buffer.WriteByte('\n')
		}
		i++
		buffer.WriteString(frame.Function)
		buffer.WriteByte('\n')
		buffer.WriteByte('\t')
		buffer.WriteString(frame.File)
		buffer.WriteByte(':')
		buffer.WriteString(strconv.Itoa(frame.Line))
	}
	return CapturedStacktrace(buffer.String())
}
