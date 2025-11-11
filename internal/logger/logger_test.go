package logger

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected LogLevel
	}{
		{"Set to None", LogLevelNone, LogLevelNone},
		{"Set to Info", LogLevelInfo, LogLevelInfo},
		{"Set to Debug", LogLevelDebug, LogLevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			if currentLevel != tt.expected {
				t.Errorf("SetLevel(%v) = %v, want %v", tt.level, currentLevel, tt.expected)
			}
		})
	}
}

func TestSetVerbose(t *testing.T) {
	SetVerbose()
	if currentLevel != LogLevelInfo {
		t.Errorf("SetVerbose() set level to %v, want %v", currentLevel, LogLevelInfo)
	}
}

func TestSetDebug(t *testing.T) {
	SetDebug()
	if currentLevel != LogLevelDebug {
		t.Errorf("SetDebug() set level to %v, want %v", currentLevel, LogLevelDebug)
	}
}

func TestIsVerbose(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"None level is not verbose", LogLevelNone, false},
		{"Info level is verbose", LogLevelInfo, true},
		{"Debug level is verbose", LogLevelDebug, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			result := IsVerbose()
			if result != tt.expected {
				t.Errorf("IsVerbose() with level %v = %v, want %v", tt.level, result, tt.expected)
			}
		})
	}
}

func TestIsDebug(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"None level is not debug", LogLevelNone, false},
		{"Info level is not debug", LogLevelInfo, false},
		{"Debug level is debug", LogLevelDebug, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			result := IsDebug()
			if result != tt.expected {
				t.Errorf("IsDebug() with level %v = %v, want %v", tt.level, result, tt.expected)
			}
		})
	}
}

func TestInfo_LogsAtInfoLevel(t *testing.T) {
	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	infoLogger = log.New(&buf, "", 0)

	// Set to Info level
	SetLevel(LogLevelInfo)

	// Log a message
	Info("Test info message: %s", "hello")

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr
	r.Close()

	// Check output
	output := buf.String()
	if !strings.Contains(output, "[INFO] Test info message: hello") {
		t.Errorf("Info() output = %q, want to contain '[INFO] Test info message: hello'", output)
	}
}

func TestInfo_DoesNotLogAtNoneLevel(t *testing.T) {
	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	infoLogger = log.New(&buf, "", 0)

	// Set to None level
	SetLevel(LogLevelNone)

	// Log a message
	Info("This should not appear")

	// Check output is empty
	output := buf.String()
	if output != "" {
		t.Errorf("Info() at None level output = %q, want empty string", output)
	}
}

func TestDebug_LogsAtDebugLevel(t *testing.T) {
	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	debugLogger = log.New(&buf, "", 0)

	// Set to Debug level
	SetLevel(LogLevelDebug)

	// Log a message
	Debug("Test debug message: %s", "world")

	// Check output
	output := buf.String()
	if !strings.Contains(output, "[DEBUG] Test debug message: world") {
		t.Errorf("Debug() output = %q, want to contain '[DEBUG] Test debug message: world'", output)
	}
}

func TestDebug_DoesNotLogAtInfoLevel(t *testing.T) {
	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	debugLogger = log.New(&buf, "", 0)

	// Set to Info level
	SetLevel(LogLevelInfo)

	// Log a message
	Debug("This should not appear")

	// Check output is empty
	output := buf.String()
	if output != "" {
		t.Errorf("Debug() at Info level output = %q, want empty string", output)
	}
}

func TestDebug_DoesNotLogAtNoneLevel(t *testing.T) {
	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	debugLogger = log.New(&buf, "", 0)

	// Set to None level
	SetLevel(LogLevelNone)

	// Log a message
	Debug("This should not appear")

	// Check output is empty
	output := buf.String()
	if output != "" {
		t.Errorf("Debug() at None level output = %q, want empty string", output)
	}
}

func TestInfo_WithMultipleArguments(t *testing.T) {
	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	infoLogger = log.New(&buf, "", 0)

	// Set to Info level
	SetLevel(LogLevelInfo)

	// Log a message with multiple format args
	Info("Request completed with status: %d in %dms", 200, 150)

	// Check output
	output := buf.String()
	expected := "[INFO] Request completed with status: 200 in 150ms"
	if !strings.Contains(output, expected) {
		t.Errorf("Info() output = %q, want to contain %q", output, expected)
	}
}

func TestDebug_WithMultipleArguments(t *testing.T) {
	// Set up a buffer to capture logger output
	var buf bytes.Buffer
	debugLogger = log.New(&buf, "", 0)

	// Set to Debug level
	SetLevel(LogLevelDebug)

	// Log a message with multiple format args
	Debug("Checking %d required fields: %v", 3, []string{"login", "id", "node_id"})

	// Check output
	output := buf.String()
	if !strings.Contains(output, "[DEBUG] Checking 3 required fields:") {
		t.Errorf("Debug() output = %q, want to contain field count", output)
	}
}

// Cleanup function to reset logger state after tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Cleanup: reset loggers to default state
	infoLogger = log.New(os.Stderr, "", 0)
	debugLogger = log.New(os.Stderr, "", 0)
	currentLevel = LogLevelNone

	os.Exit(code)
}
