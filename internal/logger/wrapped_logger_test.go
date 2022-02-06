package logger_test

import (
	"bytes"
	"log"
	"testing"

	logger "github.com/da-moon/go-template/internal/logger"
	assert "github.com/stretchr/testify/assert"
)

func TestWrappedlogger(t *testing.T) {
	testCases := []struct {
		line     string
		expected string
		level    logger.LogLevel
	}{
		{"foo\n", "[ WARN ] foo\n", logger.WarnLevel},
		{" bar\n", "[ ERROR ] bar\n", logger.ErrorLevel},
		{"baz\n", "[ DEBUG ] baz\n", logger.DebugLevel},
		{" foo\n", "[ TRACE ] foo\n", logger.TraceLevel},
		{"bar\n", "[ INFO  ] bar\n", logger.InfoLevel},
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	l := log.New(buf, "", 0)
	ll := logger.NewWrappedLogger(l)
	for _, testCase := range testCases {
		switch testCase.level {
		case logger.ErrorLevel:
			ll.Error(testCase.line)
		case logger.WarnLevel:
			ll.Warn(testCase.line)
		case logger.TraceLevel:
			ll.Trace(testCase.line)
		case logger.DebugLevel:
			ll.Debug(testCase.line)
		case logger.InfoLevel:
			ll.Info(testCase.line)
		default:
			t.Fatal("unacceptable level")
		}
		actual := buf.String()
		assert.Equal(t, testCase.expected, actual)
		buf.Reset()
	}
}
