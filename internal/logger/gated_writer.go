package logger

import (
	"io"
	"sync"
)

// GatedWriter  is an io.Writer implementation that buffers all of its
// data into an internal buffer until it is told to let data through.
// it's used to log a daemon's stdout/stderr in a more orderly fashion
type GatedWriter struct {
	writer io.Writer
	buf    [][]byte
	flush  bool
	lock   sync.RWMutex
	// lock semaphore.Semaphore
}

// NewGatedWriter returns a new gated writer
func NewGatedWriter(writer io.Writer) *GatedWriter {
	return &GatedWriter{
		writer: writer,
		// lock:   semaphore.NewBinarySemaphore(),
	}
}

// Flush ...
func (w *GatedWriter) Flush() {
	w.lock.Lock()
	// w.lock.Wait()
	w.flush = true
	w.lock.Unlock()
	// w.lock.Signal()
	for _, p := range w.buf {
		w.Write(p) //nolint:errcheck
	}
	w.buf = nil
}

// Write ...
func (w *GatedWriter) Write(p []byte) (n int, err error) {
	w.lock.RLock()
	// w.lock.Wait()
	defer w.lock.RUnlock()
	if w.flush {
		// w.lock.Signal()
		return w.writer.Write(p)
	}
	p2 := make([]byte, len(p))
	copy(p2, p)
	w.buf = append(w.buf, p2)
	// w.lock.Signal()
	return len(p), nil
}
