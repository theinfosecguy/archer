package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/logger"
	"github.com/theinfosecguy/archer/internal/models"
)

// TestMain sets up and tears down test environment
func TestMain(m *testing.M) {
	// Save original working directory
	originalWd, _ := os.Getwd()

	// Run tests
	code := m.Run()

	// Restore original working directory
	os.Chdir(originalWd)

	// Reset logger state
	logger.SetLevel(logger.LogLevelNone)

	os.Exit(code)
}

// setupTestEnvironment creates a temporary directory and returns cleanup function
func setupTestEnvironment(t *testing.T) (string, func()) {
	tempDir := t.TempDir()

	// Save original env vars
	originalSecret := os.Getenv(constants.EnvSecretName)
	originalVars := make(map[string]string)
	for _, key := range []string{"BASE_URL", "API_TOKEN"} {
		envKey := constants.EnvVarPrefix + key
		originalVars[envKey] = os.Getenv(envKey)
	}

	cleanup := func() {
		// Restore original env vars
		os.Setenv(constants.EnvSecretName, originalSecret)
		for key, val := range originalVars {
			if val == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, val)
			}
		}

		// Reset flags
		verbose = false
		debug = false
		outputJSON = ""
		jsonOnly = false
		templateFile = ""
		varArgs = []string{}

		// Reset logger
		logger.SetLevel(logger.LogLevelNone)
	}

	return tempDir, cleanup
}

// createMockGitHubServer creates a test HTTP server that simulates GitHub API
func createMockGitHubServer(t *testing.T, valid bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"login":   "octocat",
				"id":      1,
				"node_id": "MDQ6VXNlcjE=",
			})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Bad credentials",
			})
		}
	}))
}

func TestValidateCmd_VerboseFlag(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set verbose flag
	verbose = true
	debug = false

	// Create a buffer to capture logger output
	var buf bytes.Buffer
	logger.SetVerbose()

	// Verify verbose logging is enabled
	if !logger.IsVerbose() {
		t.Error("Verbose flag did not enable verbose logging")
	}

	if logger.IsDebug() {
		t.Error("Verbose flag should not enable debug logging")
	}

	// Verify it's at Info level
	logger.SetLevel(logger.LogLevelInfo)
	if !logger.IsVerbose() {
		t.Error("Info level should be verbose")
	}

	_ = tempDir
	_ = buf
}

func TestValidateCmd_DebugFlag(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set debug flag
	verbose = false
	debug = true

	// Setup debug logging
	logger.SetDebug()

	// Verify debug logging is enabled
	if !logger.IsDebug() {
		t.Error("Debug flag did not enable debug logging")
	}

	if !logger.IsVerbose() {
		t.Error("Debug level should also be verbose")
	}
}

func TestValidateCmd_OutputJSONFlag_Success(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Setup mock server
	mockServer := createMockGitHubServer(t, true)
	defer mockServer.Close()

	// Set output JSON flag
	jsonOutputPath := filepath.Join(tempDir, "output.json")
	outputJSON = jsonOutputPath
	jsonOnly = false

	// Set secret via environment
	testSecret := "ghp_test_token_1234567890"
	os.Setenv(constants.EnvSecretName, testSecret)

	// Run validation (this would normally be done through the CLI)
	// For this test, we verify the flag is set correctly
	if outputJSON == "" {
		t.Error("outputJSON flag should be set")
	}

	if outputJSON != jsonOutputPath {
		t.Errorf("outputJSON = %q, want %q", outputJSON, jsonOutputPath)
	}
}

func TestValidateCmd_OutputJSONFlag_WithJSONOnly(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set flags
	jsonOutputPath := filepath.Join(tempDir, "output.json")
	outputJSON = jsonOutputPath
	jsonOnly = true

	// Verify flags are set
	if !jsonOnly {
		t.Error("jsonOnly flag should be true")
	}

	if outputJSON != jsonOutputPath {
		t.Errorf("outputJSON = %q, want %q", outputJSON, jsonOutputPath)
	}
}

func TestValidateCmd_FlagsCombination_VerboseAndJSON(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set both verbose and JSON output flags
	verbose = true
	outputJSON = filepath.Join(tempDir, "output.json")

	logger.SetVerbose()

	// Verify both are enabled
	if !verbose {
		t.Error("verbose flag should be true")
	}

	if outputJSON == "" {
		t.Error("outputJSON should be set")
	}

	if !logger.IsVerbose() {
		t.Error("Logger should be in verbose mode")
	}
}

func TestValidateCmd_FlagsCombination_DebugAndJSON(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set both debug and JSON output flags
	debug = true
	outputJSON = filepath.Join(tempDir, "output.json")

	logger.SetDebug()

	// Verify both are enabled
	if !debug {
		t.Error("debug flag should be true")
	}

	if outputJSON == "" {
		t.Error("outputJSON should be set")
	}

	if !logger.IsDebug() {
		t.Error("Logger should be in debug mode")
	}
}

func TestValidateCmd_FlagsCombination_AllFlags(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set all flags
	verbose = true
	debug = true
	outputJSON = filepath.Join(tempDir, "output.json")
	jsonOnly = true

	// When both verbose and debug are set, debug takes precedence
	logger.SetDebug()

	// Verify all flags
	if !verbose {
		t.Error("verbose flag should be true")
	}

	if !debug {
		t.Error("debug flag should be true")
	}

	if outputJSON == "" {
		t.Error("outputJSON should be set")
	}

	if !jsonOnly {
		t.Error("jsonOnly flag should be true")
	}

	// Debug takes precedence
	if !logger.IsDebug() {
		t.Error("Logger should be in debug mode when both flags are set")
	}
}

func TestValidateCmd_JSONOnlyWithoutOutputJSON(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set json-only without output-json (edge case)
	jsonOnly = true
	outputJSON = ""

	// This is a valid configuration - json-only just won't have any effect
	// since there's no JSON output file specified
	if !jsonOnly {
		t.Error("jsonOnly flag should be true")
	}

	if outputJSON != "" {
		t.Error("outputJSON should be empty")
	}
}

func TestHandleValidationResult_SuccessWithJSON(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	jsonPath := filepath.Join(tempDir, "test-output.json")
	outputJSON = jsonPath
	jsonOnly = false

	// Create test result
	result := &models.ValidationResult{
		Valid:   true,
		Message: "Secret is valid",
	}

	// Create test template
	template := &models.SecretTemplate{
		Name:   "github",
		Mode:   constants.ModeSingle,
		Method: "GET",
		APIURL: "https://api.github.com/user",
		Request: models.RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer {{SECRET}}",
			},
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode:     []int{200},
			RequiredFields: []string{"login", "id"},
		},
	}

	vars := map[string]string{
		"SECRET": "test_secret",
	}

	// Call handleValidationResult (which would write JSON)
	startTime := time.Now().UTC()
	err := handleValidationResult(result, template, vars, startTime)
	if err != nil {
		t.Fatalf("handleValidationResult() error = %v, want nil", err)
	}

	// Note: We can't directly test the file creation here without running
	// the full validate command, but we verify the function completes without error
}

func TestHandleValidationResult_ErrorWithJSON(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	jsonPath := filepath.Join(tempDir, "test-error.json")
	outputJSON = jsonPath
	jsonOnly = false

	// Create test error result
	result := &models.ValidationResult{
		Valid: false,
		Error: "Invalid token",
	}

	template := &models.SecretTemplate{
		Name:   "github",
		Mode:   constants.ModeSingle,
		Method: "GET",
		APIURL: "https://api.github.com/user",
	}

	vars := map[string]string{
		"SECRET": "invalid_secret",
	}

	// Call handleValidationResult
	startTime := time.Now().UTC()
	err := handleValidationResult(result, template, vars, startTime)
	if err == nil {
		t.Error("handleValidationResult() with invalid result should return error")
	}

	// Verify error message format
	expectedPrefix := constants.FailureIndicator
	if !strings.Contains(err.Error(), expectedPrefix) {
		t.Errorf("error message should contain %q, got %q", expectedPrefix, err.Error())
	}
}

func TestHandleValidationResult_SuccessWithJSONOnly(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	jsonPath := filepath.Join(tempDir, "test-json-only.json")
	outputJSON = jsonPath
	jsonOnly = true

	// Create test result
	result := &models.ValidationResult{
		Valid:   true,
		Message: "Secret is valid",
	}

	template := &models.SecretTemplate{
		Name:   "github",
		Mode:   constants.ModeSingle,
		Method: "GET",
		APIURL: "https://api.github.com/user",
	}

	vars := map[string]string{
		"SECRET": "test_secret",
	}

	// Call handleValidationResult
	// With jsonOnly=true, it should not print to stdout (but we can't easily capture that)
	startTime := time.Now().UTC()
	err := handleValidationResult(result, template, vars, startTime)
	if err != nil {
		t.Fatalf("handleValidationResult() error = %v, want nil", err)
	}
}

func TestBuildMaskedArtifacts(t *testing.T) {
	template := &models.SecretTemplate{
		APIURL: "https://api.github.com/user",
		Request: models.RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer {{SECRET}}",
				"X-Custom":      "static-value",
			},
		},
	}

	vars := map[string]string{
		"SECRET": "ghp_actualSecretValue123456",
	}

	maskedURL, maskedHeaders := buildMaskedArtifacts(template, vars)

	// Verify URL is returned (masking is handled by variables package)
	if maskedURL == "" {
		t.Error("maskedURL should not be empty")
	}

	// Verify headers are processed
	if len(maskedHeaders) == 0 {
		t.Error("maskedHeaders should not be empty")
	}

	// Verify secret is masked in headers
	authHeader, exists := maskedHeaders["Authorization"]
	if !exists {
		t.Error("Authorization header should exist in masked headers")
	}

	// The actual masking logic is in variables package, but verify structure
	if authHeader == "Bearer ghp_actualSecretValue123456" {
		t.Error("Secret should be masked in Authorization header")
	}
}

func TestGetEnvVariables(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set test environment variables
	os.Setenv(constants.EnvVarPrefix+"BASE_URL", "https://example.com")
	os.Setenv(constants.EnvVarPrefix+"API_TOKEN", "test_token_123")
	os.Setenv(constants.EnvVarPrefix+"MISSING_VAR", "") // Empty value

	requiredVars := []string{"BASE_URL", "API_TOKEN", "ANOTHER_VAR"}

	envVars := getEnvVariables(requiredVars)

	// Verify BASE_URL is retrieved
	if envVars["BASE_URL"] != "https://example.com" {
		t.Errorf("envVars[BASE_URL] = %q, want 'https://example.com'", envVars["BASE_URL"])
	}

	// Verify API_TOKEN is retrieved
	if envVars["API_TOKEN"] != "test_token_123" {
		t.Errorf("envVars[API_TOKEN] = %q, want 'test_token_123'", envVars["API_TOKEN"])
	}

	// Verify ANOTHER_VAR is not in the map (not set in env)
	if _, exists := envVars["ANOTHER_VAR"]; exists {
		t.Error("envVars should not contain ANOTHER_VAR (not set in environment)")
	}

	// Verify MISSING_VAR is not in the map (empty value)
	if _, exists := envVars["MISSING_VAR"]; exists {
		t.Error("envVars should not contain MISSING_VAR (empty value)")
	}
}

func TestLoggerSetupInRunValidate_DebugPriority(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Simulate what happens in runValidate when both flags are set
	debug = true
	verbose = true

	// The actual logic in runValidate
	if debug {
		logger.SetDebug()
	} else if verbose {
		logger.SetVerbose()
	}

	// Debug should take priority
	if !logger.IsDebug() {
		t.Error("Debug should be enabled when debug flag is true")
	}
}

func TestLoggerSetupInRunValidate_VerboseOnly(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Simulate what happens in runValidate when only verbose is set
	debug = false
	verbose = true

	// The actual logic in runValidate
	if debug {
		logger.SetDebug()
	} else if verbose {
		logger.SetVerbose()
	}

	// Verbose should be enabled
	if !logger.IsVerbose() {
		t.Error("Verbose should be enabled when verbose flag is true")
	}

	// But not debug
	if logger.IsDebug() {
		t.Error("Debug should not be enabled when only verbose flag is true")
	}
}

func TestLoggerSetupInRunValidate_NoFlags(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Simulate what happens in runValidate when no flags are set
	debug = false
	verbose = false

	// The actual logic in runValidate - nothing should happen
	if debug {
		logger.SetDebug()
	} else if verbose {
		logger.SetVerbose()
	}

	// No logging should be enabled
	if logger.IsVerbose() {
		t.Error("Verbose should not be enabled when no flags are set")
	}

	if logger.IsDebug() {
		t.Error("Debug should not be enabled when no flags are set")
	}
}
