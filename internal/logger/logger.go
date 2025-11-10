package logger

import (
	"fmt"
	"log"
	"os"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// LogLevelNone - no logging
	LogLevelNone LogLevel = iota
	// LogLevelInfo - info level logging (verbose mode)
	LogLevelInfo
	// LogLevelDebug - debug level logging (debug mode)
	LogLevelDebug
)

var (
	currentLevel LogLevel = LogLevelNone
	infoLogger   *log.Logger
	debugLogger  *log.Logger
)

func init() {
	// Initialize loggers with no prefix, output to stderr
	infoLogger = log.New(os.Stderr, "", 0)
	debugLogger = log.New(os.Stderr, "", 0)
}

// SetLevel sets the current logging level
func SetLevel(level LogLevel) {
	currentLevel = level
}

// SetVerbose enables verbose (INFO) logging
func SetVerbose() {
	currentLevel = LogLevelInfo
}

// SetDebug enables debug logging
func SetDebug() {
	currentLevel = LogLevelDebug
}

// IsVerbose returns true if verbose logging is enabled
func IsVerbose() bool {
	return currentLevel >= LogLevelInfo
}

// IsDebug returns true if debug logging is enabled
func IsDebug() bool {
	return currentLevel >= LogLevelDebug
}

// Info logs an info-level message (shown in verbose mode)
func Info(format string, args ...interface{}) {
	if currentLevel >= LogLevelInfo {
		msg := fmt.Sprintf(format, args...)
		infoLogger.Printf("[INFO] %s", msg)
	}
}

// Debug logs a debug-level message (shown in debug mode)
func Debug(format string, args ...interface{}) {
	if currentLevel >= LogLevelDebug {
		msg := fmt.Sprintf(format, args...)
		debugLogger.Printf("[DEBUG] %s", msg)
	}
}
