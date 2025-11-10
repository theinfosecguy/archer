package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/models"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.restyClient == nil {
		t.Fatal("restyClient is nil")
	}
}

func TestExecuteRequest_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header, got %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user_id": "123", "valid": true}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${SECRET}",
			},
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode:     []int{200},
			RequiredFields: []string{"$.user_id", "$.valid"},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
		},
	}

	vars := map[string]string{
		"SECRET": "test-token",
	}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid result, got invalid: %s", result.Error)
	}

	if result.Message != constants.SecretValid {
		t.Errorf("Expected message '%s', got '%s'", constants.SecretValid, result.Message)
	}
}

func TestExecuteRequest_InvalidStatusCode(t *testing.T) {
	// Create test server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
			ErrorMessages: map[int]string{
				401: "Invalid API key",
			},
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result, got valid")
	}

	if result.Error != "Invalid API key" {
		t.Errorf("Expected custom error message 'Invalid API key', got '%s'", result.Error)
	}
}

func TestExecuteRequest_MissingRequiredField(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user_id": "123"}`)) // Missing "valid" field
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode:     []int{200},
			RequiredFields: []string{"$.user_id", "$.valid"},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result due to missing field, got valid")
	}

	if result.Error == "" {
		t.Error("Expected error message about missing field")
	}
}

func TestExecuteRequest_WithJSONBody(t *testing.T) {
	// Create test server
	receivedBody := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		receivedBody = string(buf[:n])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "POST",
		Request: models.RequestConfig{
			JSONData: map[string]any{
				"api_key": "${API_KEY}",
				"secret":  "${API_SECRET}",
			},
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
		},
	}

	vars := map[string]string{
		"API_KEY":    "test-key",
		"API_SECRET": "test-secret",
	}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid result, got invalid: %s", result.Error)
	}

	// Verify body was sent
	if receivedBody == "" {
		t.Error("Expected request body to be sent")
	}
}

func TestExecuteRequest_WithQueryParams(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query params
		if r.URL.Query().Get("api_key") != "test-key" {
			t.Errorf("Expected api_key=test-key, got %s", r.URL.Query().Get("api_key"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			QueryParams: map[string]string{
				"api_key": "${SECRET}",
			},
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
		},
	}

	vars := map[string]string{
		"SECRET": "test-key",
	}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid result, got invalid: %s", result.Error)
	}
}

func TestExecuteRequest_InvalidJSON(t *testing.T) {
	// Create test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode:     []int{200},
			RequiredFields: []string{"$.user_id"},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result due to invalid JSON, got valid")
	}

	if result.Error != constants.InvalidJSONResponse {
		t.Errorf("Expected error '%s', got '%s'", constants.InvalidJSONResponse, result.Error)
	}
}

func TestCheckResponse_MultipleStatusCodes(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated) // 201
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "POST",
		Request: models.RequestConfig{
			Timeout: 10,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200, 201, 202}, // Accept multiple status codes
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0,
			RetryDelay: 0,
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid result for 201 status, got invalid: %s", result.Error)
	}
}

func TestExecuteRequest_WithRetries(t *testing.T) {
	attempts := 0
	// Create test server that fails first 2 attempts, succeeds on 3rd
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 5,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 3,
			RetryDelay: 0, // No delay for faster test
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid result after retries, got invalid: %s", result.Error)
	}

	// Verify it actually retried (3 attempts total: initial + 2 retries)
	if attempts != 3 {
		t.Errorf("Expected 3 attempts (with 2 retries), got %d", attempts)
	}
}

func TestExecuteRequest_RetriesExhausted(t *testing.T) {
	attempts := 0
	// Create test server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 5,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 2,
			RetryDelay: 0, // No delay for faster test
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result after exhausting retries, got valid")
	}

	// Verify it retried the correct number of times (initial + 2 retries = 3 total)
	if attempts != 3 {
		t.Errorf("Expected 3 attempts (initial + 2 retries), got %d", attempts)
	}
}

func TestExecuteRequest_TimeoutError(t *testing.T) {
	// Create test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 1, // 1 second timeout
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0, // No retries for timeout test
			RetryDelay: 0,
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result due to timeout, got valid")
	}

	if result.Error != constants.RequestTimeout {
		t.Errorf("Expected error '%s', got '%s'", constants.RequestTimeout, result.Error)
	}
}

func TestExecuteRequest_NoRetries(t *testing.T) {
	attempts := 0
	// Create test server that fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient()
	template := &models.SecretTemplate{
		APIURL: server.URL,
		Method: "GET",
		Request: models.RequestConfig{
			Timeout: 5,
		},
		SuccessCriteria: models.SuccessCriteria{
			StatusCode: []int{200},
		},
		ErrorHandling: models.ErrorHandling{
			MaxRetries: 0, // No retries
			RetryDelay: 0,
		},
	}

	vars := map[string]string{}

	result, err := client.ExecuteRequest(template, vars)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result, got valid")
	}

	// Verify it only attempted once (no retries)
	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retries), got %d", attempts)
	}
}
