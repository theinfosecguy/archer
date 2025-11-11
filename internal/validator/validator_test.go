package validator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSecretValidator(t *testing.T) {
	templatesDir := "testdata/templates"
	validator := NewSecretValidator(templatesDir)

	if validator == nil {
		t.Fatal("NewSecretValidator() returned nil")
	}

	if validator.TemplateLoader == nil {
		t.Error("TemplateLoader is nil")
	}

	if validator.HTTPClient == nil {
		t.Error("HTTPClient is nil")
	}
}

func TestNewSecretValidator_EmptyDirectory(t *testing.T) {
	validator := NewSecretValidator("")

	if validator == nil {
		t.Fatal("NewSecretValidator('') returned nil")
	}

	if validator.TemplateLoader.DefaultTemplatesDir != "templates" {
		t.Errorf("DefaultTemplatesDir = %q, want 'templates'", validator.TemplateLoader.DefaultTemplatesDir)
	}
}

func TestValidateSecret_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer ghp_kLm3N4pQ5rS6tU7vW8xY9zA0bC1dE2fG3hI4jK5" {
			t.Errorf("Authorization header = %q, want Bearer ghp_kLm3N4pQ5rS6tU7vW8xY9zA0bC1dE2fG3hI4jK5", r.Header.Get("Authorization"))
		}

		response := map[string]any{
			"login":   "octocat",
			"id":      1,
			"node_id": "MDQ6VXNlcjE=",
			"name":    "The Octocat",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("github")
	template.APIURL = mockServer.URL

	vars := map[string]string{"SECRET": "ghp_kLm3N4pQ5rS6tU7vW8xY9zA0bC1dE2fG3hI4jK5"}
	result, err := validator.validateWithTemplate(template, vars)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if !result.Valid {
		t.Errorf("result.Valid = false, want true. Error: %s", result.Error)
	}
}

func TestValidateSecret_InvalidToken(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad credentials",
		})
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("github")
	template.APIURL = mockServer.URL

	vars := map[string]string{"SECRET": "ghp_expired_4aB9cD2eF7gH1iJ6kL3mN8oP5qR2sT4uV"}
	result, err := validator.validateWithTemplate(template, vars)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error != "Invalid or expired GitHub token" {
		t.Errorf("result.Error = %q, want 'Invalid or expired GitHub token'", result.Error)
	}
}

func TestValidateSecret_TemplateNotFound(t *testing.T) {
	validator := NewSecretValidator("testdata/templates")
	result, err := validator.ValidateSecret("nonexistent", "ghp_nonexistent_3aB7cD9eF2gH5iJ8kL1mN4oP6qR")

	if err != nil {
		t.Fatalf("ValidateSecret() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error == "" {
		t.Error("result.Error is empty, want template not found error")
	}
}

func TestValidateSecret_MultipartTemplateRejected(t *testing.T) {
	validator := NewSecretValidator("testdata/templates")
	result, err := validator.ValidateSecret("ghost", "some_token")

	if err != nil {
		t.Fatalf("ValidateSecret() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error == "" {
		t.Error("result.Error is empty, want mode error")
	}
}

func TestValidateSecret_StatusCodeMismatch(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Forbidden",
		})
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("github")
	template.APIURL = mockServer.URL

	vars := map[string]string{"SECRET": "ghp_revoked_9zX8yW7vU6tS5rQ4pO3nM2lK1jI0hG"}
	result, err := validator.validateWithTemplate(template, vars)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error != "Token lacks required permissions" {
		t.Errorf("result.Error = %q, want 'Token lacks required permissions'", result.Error)
	}
}

func TestValidateSecret_MissingRequiredField(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"login": "octocat",
			"id":    1,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("github")
	template.APIURL = mockServer.URL

	vars := map[string]string{"SECRET": "ghp_malformed_8fE4dC2bA1zY9xW7vU5tS3rQ1"}
	result, err := validator.validateWithTemplate(template, vars)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error == "" {
		t.Error("result.Error is empty, want required field error")
	}
}

func TestValidateSecret_InvalidJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json response"))
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("github")
	template.APIURL = mockServer.URL

	vars := map[string]string{"SECRET": "ghp_malformed_8fE4dC2bA1zY9xW7vU5tS3rQ1"}
	result, err := validator.validateWithTemplate(template, vars)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error == "" {
		t.Error("result.Error is empty, want invalid JSON error")
	}
}

func TestValidateSecret_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("github")
	template.APIURL = mockServer.URL

	vars := map[string]string{"SECRET": "ghp_token"}
	result, err := validator.validateWithTemplate(template, vars)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false (timeout)")
	}

	if result.Error != "Request timeout" {
		t.Errorf("result.Error = %q, want 'Request timeout'", result.Error)
	}
}

func TestValidateSecretMultipart_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "62f8a45c3e1d9b7a4f6e2c8d5a9b3e7f" {
			t.Errorf("query param key = %q, want 62f8a45c3e1d9b7a4f6e2c8d5a9b3e7f", r.URL.Query().Get("key"))
		}

		if r.URL.Query().Get("limit") != "1" {
			t.Errorf("query param limit = %q, want 1", r.URL.Query().Get("limit"))
		}

		response := map[string]any{
			"posts": []map[string]any{
				{"id": "1", "title": "Test Post"},
			},
			"meta": map[string]any{
				"pagination": map[string]int{"total": 1},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("ghost")
	template.APIURL = mockServer.URL + "/ghost/api/content/posts/"

	variables := map[string]string{
		"BASE_URL":  mockServer.URL,
		"API_TOKEN": "62f8a45c3e1d9b7a4f6e2c8d5a9b3e7f",
	}

	result, err := validator.validateWithTemplate(template, variables)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if !result.Valid {
		t.Errorf("result.Valid = false, want true. Error: %s", result.Error)
	}
}

func TestValidateSecretMultipart_MissingVariable(t *testing.T) {
	validator := NewSecretValidator("testdata/templates")
	variables := map[string]string{
		"BASE_URL": "https://techcrunch.ghost.io",
	}

	result, err := validator.ValidateSecretMultipart("ghost", variables)

	if err != nil {
		t.Fatalf("ValidateSecretMultipart() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error == "" {
		t.Error("result.Error is empty, want missing variable error")
	}
}

func TestValidateSecretMultipart_TemplateNotFound(t *testing.T) {
	validator := NewSecretValidator("testdata/templates")
	variables := map[string]string{
		"API_KEY": "key123",
	}

	result, err := validator.ValidateSecretMultipart("nonexistent", variables)

	if err != nil {
		t.Fatalf("ValidateSecretMultipart() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}
}

func TestValidateSecretMultipart_SingleModeTemplateRejected(t *testing.T) {
	validator := NewSecretValidator("testdata/templates")
	variables := map[string]string{
		"SECRET": "ghp_wrongmode_5xY9wV2uT8sR4qP1nM7lK3jI6hG",
	}

	result, err := validator.ValidateSecretMultipart("github", variables)

	if err != nil {
		t.Fatalf("ValidateSecretMultipart() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}

	if result.Error == "" {
		t.Error("result.Error is empty, want mode error")
	}
}

func TestValidateSecretMultipart_AllVariablesMissing(t *testing.T) {
	validator := NewSecretValidator("testdata/templates")
	variables := map[string]string{}

	result, err := validator.ValidateSecretMultipart("ghost", variables)

	if err != nil {
		t.Fatalf("ValidateSecretMultipart() error = %v, want nil", err)
	}

	if result.Valid {
		t.Errorf("result.Valid = true, want false")
	}
}

func TestValidateSecretMultipart_QueryParamsInjected(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("key") == "" {
			t.Error("query param 'key' is empty")
		}

		if query.Get("limit") != "1" {
			t.Errorf("query param 'limit' = %q, want '1'", query.Get("limit"))
		}

		response := map[string]any{
			"posts": []any{},
			"meta":  map[string]any{"total": 0},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	validator := NewSecretValidator("testdata/templates")
	template, _ := validator.TemplateLoader.GetTemplate("ghost")
	template.APIURL = mockServer.URL

	variables := map[string]string{
		"BASE_URL":  mockServer.URL,
		"API_TOKEN": "8f3e9d7c6b5a4f2e1d9c8b7a6f5e4d3c",
	}

	result, err := validator.validateWithTemplate(template, variables)

	if err != nil {
		t.Fatalf("validateWithTemplate() error = %v, want nil", err)
	}

	if !result.Valid {
		t.Errorf("result.Valid = false, want true. Error: %s", result.Error)
	}
}
