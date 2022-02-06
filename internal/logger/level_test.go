package logger_test

import (
	"bytes"
	"io"
	"log"
	"testing"

	logger "github.com/da-moon/go-template/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestLevelFilter_impl(t *testing.T) {
	var _ io.Writer = logger.NewLevelFilter()
}
func TestLevelFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := logger.NewLevelFilter(
		logger.WithWriter(buf),
		logger.WithMinLevel("WARN"),
		logger.WithLevels([]string{"DEBUG", "WARN", "ERROR"}),
	)
	logger := log.New(filter, "", 0)
	logger.Print("[ WARN ] foo")
	logger.Println("[ ERROR ] bar")
	logger.Println("[ DEBUG ] baz")
	logger.Println("[ WARN ] buzz")
	actual := buf.String()
	expected := "[ WARN ] foo\n[ ERROR ] bar\n[ WARN ] buzz\n"
	assert.Equal(t, expected, actual)
}
func TestLevelFilterCheck(t *testing.T) {
	filter := logger.NewLevelFilter(
		logger.WithMinLevel("WARN"),
		logger.WithLevels([]string{"DEBUG", "WARN", "ERROR"}),
	)
	tests := []struct {
		line  string
		check bool
	}{
		{"[ WARN ] foo\n", true},
		{"[ ERROR ] bar\n", true},
		{"[ DEBUG ] baz\n", false},
		{"[ WARN ] buzz\n", true},
	}
	for _, tt := range tests {
		actual := filter.Check([]byte(tt.line))
		assert.Equal(t, tt.check, actual)
	}
}
func TestLevelFilter_SetMinLevel(t *testing.T) {
	filter := logger.NewLevelFilter(
		logger.WithMinLevel("ERROR"),
		logger.WithLevels([]string{"DEBUG", "WARN", "ERROR"}),
	)
	tests := []struct {
		line        string
		checkBefore bool
		checkAfter  bool
	}{
		{"[ WARN ] foo\n", false, true},
		{"[ ERROR ] bar\n", true, true},
		{"[ DEBUG ] baz\n", false, false},
		{"[ WARN ] buzz\n", false, true},
	}
	for _, tt := range tests {
		actual := filter.Check([]byte(tt.line))
		assert.Equal(t, tt.checkBefore, actual, "expected=%v actual=%v", tt.checkBefore, actual)
	}
	// Update the minimum level to WARN
	filter.SetMinLevel("WARN")
	for _, tt := range tests {
		actual := filter.Check([]byte(tt.line))
		assert.Equal(t, tt.checkAfter, actual)
	}
}
