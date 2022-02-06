package logger

import "io"

// LevelFilterOption - sets levelFilter options
type LevelFilterOption func(*levelFilter)

// LogLevel string representing log level
type LogLevel string

// default log levels
var (
	// TraceLevel ...
	TraceLevel LogLevel = "TRACE" //nolint:gochecknoglobals
	// DebugLevel ...
	DebugLevel LogLevel = "DEBUG" //nolint:gochecknoglobals
	// InfoLevel ...
	InfoLevel LogLevel = "INFO" //nolint:gochecknoglobals
	// WarnLevel ...
	WarnLevel LogLevel = "WARN" //nolint:gochecknoglobals
	// ErrorLevel ...
	ErrorLevel LogLevel = "ERROR" //nolint:gochecknoglobals
)

// DefaultLogLevels ...
var DefaultLogLevels = []LogLevel{ //nolint:gochecknoglobals
	TraceLevel,
	DebugLevel,
	InfoLevel,
	WarnLevel,
	ErrorLevel,
}

// WithWriter - sets io.writer
func WithWriter(arg io.Writer) LevelFilterOption {
	return func(s *levelFilter) {
		s.writer = arg
	}
}

// WithMinLevel ...
func WithMinLevel(arg string) LevelFilterOption {
	return func(s *levelFilter) {
		s.minLevel = LogLevel(arg)
	}
}

// WithLevels ...
func WithLevels(arg []string) LevelFilterOption {
	return func(s *levelFilter) {
		if s.levels == nil {
			s.levels = make([]LogLevel, 0)
		}
		for _, v := range arg {
			s.levels = append(s.levels, LogLevel(v))
		}
	}
}
