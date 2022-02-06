package logger_test

import (
	"bytes"
	"io"
	"testing"

	logger "github.com/da-moon/go-template/internal/logger"
	"github.com/stretchr/testify/assert"
)

var _ io.Writer = &logger.GatedWriter{}

// TestGatedWriter ...
func TestGatedWriter(t *testing.T) {
	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	w := logger.NewGatedWriter(buf)
	_, err := w.Write([]byte("foo\n"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("bar\n"))
	assert.NoError(t, err)
	if buf.String() != "" {
		t.Fatalf("bad: %s", buf.String())
	}
	w.Flush()
	if buf.String() != "foo\nbar\n" {
		t.Fatalf("bad: %s", buf.String())
	}
	_, err = w.Write([]byte("baz\n"))
	assert.NoError(t, err)
	if buf.String() != "foo\nbar\nbaz\n" {
		t.Fatalf("bad: %s", buf.String())
	}
}
