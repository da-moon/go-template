package logger_test

import (
	"testing"

	logger "github.com/da-moon/go-template/internal/logger"
	assert "github.com/stretchr/testify/assert"
)

type MockLogHandler struct {
	logs []string
}

func (m *MockLogHandler) HandleLog(l string) {
	m.logs = append(m.logs, l)
}
func TestLogWriter(t *testing.T) {
	h := &MockLogHandler{}
	w := logger.NewLogWriter(4)
	_, err := w.Write([]byte("one"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("two"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("three"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("four"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("five"))
	assert.NoError(t, err)
	w.RegisterHandler(h)
	_, err = w.Write([]byte("six"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("seven"))
	assert.NoError(t, err)
	w.DeregisterHandler(h)
	_, err = w.Write([]byte("eight"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("nine"))
	assert.NoError(t, err)
	out := []string{
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
	}
	for idx := range out {
		assert.Equal(t, out[idx], h.logs[idx])
	}
}
