package logger

import (
	"strings"
)

// Level is level of logging
type Level uint16

// Log-level constant
const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

// Logger interface
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	FataL(format string, args ...interface{})
	Close()
}

// Gets the corresponding string
func getLevelStr(level Level) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

// Resolve the corresponding Level
func parseLogLevel(levelstr string) Level {
	levelStr := strings.ToLower(levelstr)
	switch levelStr {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "error":
		return ErrorLevel
	case "fata":
		return FatalLevel
	default:
		return DebugLevel

	}
}
