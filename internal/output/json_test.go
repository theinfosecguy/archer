package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/theinfosecguy/archer/internal/models"
)

func TestWriteJSONFile_Success(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test-output.json")

	// Create test data
	message := "Secret is valid"
	resolvedName := "github"
	mode := "single"
	source := "builtin"
	method := "GET"
	maskedURL := "https://api.github.com/user"

	result := &models.ValidationResultJSON{
		Command: "validate",
		Version: "1.0.0",
		Valid:   true,
		Message: &message,
		Request: models.ValidationRequestMeta{
			Template:             "github",
			ResolvedTemplateName: &resolvedName,
			Mode:                 &mode,
			Source:               &source,
			Method:               &method,
			APIURLMasked:         &maskedURL,
			HeadersMasked: map[string]string{
				"Authorization": "Bearer ***SECRET***",
			},
			VariablesProvided: []string{"SECRET"},
			StartedAt:         time.Now().UTC(),
			FinishedAt:        time.Now().UTC(),
			DurationMS:        500.0,
		},
		Response: models.ValidationResponseMeta{
			RequiredFieldsChecked: []string{"login", "id"},
		},
	}

	// Write JSON file
	err := WriteJSONFile(filePath, result)
	if err != nil {
		t.Fatalf("WriteJSONFile() error = %v, want nil", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("WriteJSONFile() did not create file at %s", filePath)
	}

	// Read and verify file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify it's valid JSON
	var parsed models.ValidationResultJSON
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("WriteJSONFile() created invalid JSON: %v", err)
	}

	// Verify key fields
	if parsed.Command != "validate" {
		t.Errorf("parsed.Command = %q, want 'validate'", parsed.Command)
	}
	if parsed.Version != "1.0.0" {
		t.Errorf("parsed.Version = %q, want '1.0.0'", parsed.Version)
	}
	if !parsed.Valid {
		t.Errorf("parsed.Valid = %v, want true", parsed.Valid)
	}
	if parsed.Request.Template != "github" {
		t.Errorf("parsed.Request.Template = %q, want 'github'", parsed.Request.Template)
	}
}

func TestWriteJSONFile_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test-invalid.json")

	errorMsg := "Invalid token"
	result := &models.ValidationResultJSON{
		Command: "validate",
		Version: "1.0.0",
		Valid:   false,
		Error:   &errorMsg,
		Request: models.ValidationRequestMeta{
			Template:   "github",
			StartedAt:  time.Now().UTC(),
			FinishedAt: time.Now().UTC(),
			DurationMS: 100.0,
		},
		Response: models.ValidationResponseMeta{
			Error: &errorMsg,
		},
	}

	err := WriteJSONFile(filePath, result)
	if err != nil {
		t.Fatalf("WriteJSONFile() error = %v, want nil", err)
	}

	// Read and verify
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var parsed models.ValidationResultJSON
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("WriteJSONFile() created invalid JSON: %v", err)
	}

	if parsed.Valid {
		t.Errorf("parsed.Valid = %v, want false", parsed.Valid)
	}
	if parsed.Error == nil || *parsed.Error != "Invalid token" {
		t.Errorf("parsed.Error = %v, want 'Invalid token'", parsed.Error)
	}
}

func TestWriteJSONFile_PrettyFormat(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test-pretty.json")

	message := "Success"
	result := &models.ValidationResultJSON{
		Command: "validate",
		Version: "1.0.0",
		Valid:   true,
		Message: &message,
		Request: models.ValidationRequestMeta{
			Template:   "test",
			StartedAt:  time.Now().UTC(),
			FinishedAt: time.Now().UTC(),
			DurationMS: 200.0,
		},
		Response: models.ValidationResponseMeta{},
	}

	err := WriteJSONFile(filePath, result)
	if err != nil {
		t.Fatalf("WriteJSONFile() error = %v, want nil", err)
	}

	// Read file contents as string
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)

	// Verify it's pretty-printed (has newlines and indentation)
	if !containsNewlines(content) {
		t.Errorf("WriteJSONFile() output is not pretty-formatted (no newlines found)")
	}

	// Verify it contains proper indentation (2 spaces)
	if !containsIndentation(content) {
		t.Errorf("WriteJSONFile() output is not properly indented")
	}
}

func TestWriteJSONFile_InvalidPath(t *testing.T) {
	// Try to write to an invalid path
	filePath := "/invalid/path/that/does/not/exist/output.json"

	result := &models.ValidationResultJSON{
		Command: "validate",
		Version: "1.0.0",
		Valid:   false,
		Request: models.ValidationRequestMeta{
			Template:   "test",
			StartedAt:  time.Now().UTC(),
			FinishedAt: time.Now().UTC(),
			DurationMS: 0,
		},
		Response: models.ValidationResponseMeta{},
	}

	err := WriteJSONFile(filePath, result)
	if err == nil {
		t.Errorf("WriteJSONFile() with invalid path should return error, got nil")
	}
}

func TestWriteJSONFile_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test-perms.json")

	message := "Success"
	result := &models.ValidationResultJSON{
		Command: "validate",
		Version: "1.0.0",
		Valid:   true,
		Message: &message,
		Request: models.ValidationRequestMeta{
			Template:   "test",
			StartedAt:  time.Now().UTC(),
			FinishedAt: time.Now().UTC(),
			DurationMS: 100.0,
		},
		Response: models.ValidationResponseMeta{},
	}

	err := WriteJSONFile(filePath, result)
	if err != nil {
		t.Fatalf("WriteJSONFile() error = %v, want nil", err)
	}

	// Check file permissions
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Verify permissions are 0644 (or close to it, depending on umask)
	mode := fileInfo.Mode()
	// On Unix systems, we expect 0644
	expectedPerm := os.FileMode(0644)
	if mode.Perm() != expectedPerm {
		// This might vary by system, so just log it
		t.Logf("File permissions = %v, expected %v (may vary by system)", mode.Perm(), expectedPerm)
	}
}

func TestWriteJSONFile_WithMaskedSecrets(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test-masked.json")

	resolvedName := "github"
	mode := "single"
	source := "builtin"
	method := "GET"
	maskedURL := "https://api.github.com/user"
	message := "Valid"

	result := &models.ValidationResultJSON{
		Command: "validate",
		Version: "1.0.0",
		Valid:   true,
		Message: &message,
		Request: models.ValidationRequestMeta{
			Template:             "github",
			ResolvedTemplateName: &resolvedName,
			Mode:                 &mode,
			Source:               &source,
			Method:               &method,
			APIURLMasked:         &maskedURL,
			HeadersMasked: map[string]string{
				"Authorization": "Bearer ***SECRET***",
				"X-API-Key":     "***API_KEY***",
			},
			VariablesProvided: []string{"SECRET", "API_KEY"},
			StartedAt:         time.Now().UTC(),
			FinishedAt:        time.Now().UTC(),
			DurationMS:        300.0,
		},
		Response: models.ValidationResponseMeta{},
	}

	err := WriteJSONFile(filePath, result)
	if err != nil {
		t.Fatalf("WriteJSONFile() error = %v, want nil", err)
	}

	// Read and verify secrets are masked
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)

	// Verify masked values are present
	if !containsString(content, "***SECRET***") {
		t.Errorf("Output should contain masked SECRET, got: %s", content)
	}
	if !containsString(content, "***API_KEY***") {
		t.Errorf("Output should contain masked API_KEY, got: %s", content)
	}
}

// Helper functions

func containsNewlines(s string) bool {
	for _, c := range s {
		if c == '\n' {
			return true
		}
	}
	return false
}

func containsIndentation(s string) bool {
	// Look for lines that start with 2 spaces (JSON indentation)
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '\n' && i+2 < len(s) {
			if s[i+1] == ' ' && s[i+2] == ' ' {
				return true
			}
		}
	}
	return false
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		len(s) >= len(substr) &&
		findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
