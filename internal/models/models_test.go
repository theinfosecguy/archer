package models

import (
	"testing"
)

func TestRequestConfig_Validate_BothDataAndJSONData(t *testing.T) {
	data := "some data"
	config := RequestConfig{
		Data:     &data,
		JSONData: map[string]any{"key": "value"},
	}

	err := config.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want mutual exclusion error")
	}
}

func TestRequestConfig_Validate_OnlyData(t *testing.T) {
	data := "grant_type=client_credentials&client_id=abc123"
	config := RequestConfig{
		Data:     &data,
		JSONData: nil,
	}

	err := config.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestRequestConfig_Validate_OnlyJSONData(t *testing.T) {
	config := RequestConfig{
		Data: nil,
		JSONData: map[string]any{
			"username": "admin",
			"password": "secret",
		},
	}

	err := config.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestRequestConfig_Validate_NeitherDataNorJSONData(t *testing.T) {
	config := RequestConfig{
		Data:     nil,
		JSONData: nil,
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
	}

	err := config.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestSecretTemplate_SetDefaults_EmptyMethod(t *testing.T) {
	template := SecretTemplate{
		Name: "test-template",
	}

	template.SetDefaults()

	if template.Method != "GET" {
		t.Errorf("Method = %q, want 'GET'", template.Method)
	}
}

func TestSecretTemplate_SetDefaults_EmptyMode(t *testing.T) {
	template := SecretTemplate{
		Name: "test-template",
	}

	template.SetDefaults()

	if template.Mode != "single" {
		t.Errorf("Mode = %q, want 'single'", template.Mode)
	}
}

func TestSecretTemplate_SetDefaults_ZeroTimeout(t *testing.T) {
	template := SecretTemplate{
		Name: "test-template",
	}

	template.SetDefaults()

	if template.Request.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30", template.Request.Timeout)
	}
}

func TestSecretTemplate_SetDefaults_ZeroMaxRetries(t *testing.T) {
	template := SecretTemplate{
		Name: "test-template",
	}

	template.SetDefaults()

	if template.ErrorHandling.MaxRetries != 0 {
		t.Errorf("MaxRetries = %d, want 0", template.ErrorHandling.MaxRetries)
	}
}

func TestSecretTemplate_SetDefaults_ZeroRetryDelay(t *testing.T) {
	template := SecretTemplate{
		Name: "test-template",
	}

	template.SetDefaults()

	if template.ErrorHandling.RetryDelay != 0 {
		t.Errorf("RetryDelay = %d, want 0", template.ErrorHandling.RetryDelay)
	}
}

func TestSecretTemplate_SetDefaults_PreservesExistingValues(t *testing.T) {
	template := SecretTemplate{
		Name:   "test-template",
		Method: "POST",
		Mode:   "multipart",
		Request: RequestConfig{
			Timeout: 60,
		},
		ErrorHandling: ErrorHandling{
			MaxRetries: 5,
			RetryDelay: 10,
		},
	}

	template.SetDefaults()

	if template.Method != "POST" {
		t.Errorf("Method = %q, want 'POST'", template.Method)
	}

	if template.Mode != "multipart" {
		t.Errorf("Mode = %q, want 'multipart'", template.Mode)
	}

	if template.Request.Timeout != 60 {
		t.Errorf("Timeout = %d, want 60", template.Request.Timeout)
	}

	if template.ErrorHandling.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", template.ErrorHandling.MaxRetries)
	}

	if template.ErrorHandling.RetryDelay != 10 {
		t.Errorf("RetryDelay = %d, want 10", template.ErrorHandling.RetryDelay)
	}
}

func TestSecretTemplate_Validate_ValidSingleMode(t *testing.T) {
	template := SecretTemplate{
		Name:   "github",
		Mode:   "single",
		APIURL: "https://api.github.com/user",
		Request: RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${SECRET}",
			},
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestSecretTemplate_Validate_ValidMultipartMode(t *testing.T) {
	template := SecretTemplate{
		Name:   "ghost",
		Mode:   "multipart",
		APIURL: "${BASE_URL}/api/posts",
		RequiredVariables: []string{
			"BASE_URL",
			"API_TOKEN",
		},
		Request: RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${API_TOKEN}",
			},
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestSecretTemplate_Validate_InvalidMode(t *testing.T) {
	template := SecretTemplate{
		Name:   "github",
		Mode:   "invalid_mode",
		APIURL: "https://api.github.com/user",
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want mode validation error")
	}
}

func TestSecretTemplate_Validate_InvalidVariableFormat(t *testing.T) {
	template := SecretTemplate{
		Name:   "test",
		Mode:   "multipart",
		APIURL: "${BASE_URL}/api",
		RequiredVariables: []string{
			"BASE_URL",
			"invalid-kebab-case",
		},
		Request: RequestConfig{
			Headers: map[string]string{
				"X-API-Key": "${invalid-kebab-case}",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want variable format error")
	}
}

func TestSecretTemplate_Validate_MultipartMissingRequiredVariables(t *testing.T) {
	template := SecretTemplate{
		Name:   "ghost",
		Mode:   "multipart",
		APIURL: "${BASE_URL}/ghost/api/content/posts/",
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want missing required_variables error")
	}
}

func TestSecretTemplate_Validate_SingleModeWithRequiredVariables(t *testing.T) {
	template := SecretTemplate{
		Name:              "openai",
		Mode:              "single",
		APIURL:            "https://api.openai.com/v1/models",
		RequiredVariables: []string{"API_KEY"},
		Request: RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${SECRET}",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want single mode cannot have required_variables error")
	}
}

func TestSecretTemplate_Validate_SingleModeWithInvalidVariable(t *testing.T) {
	template := SecretTemplate{
		Name:   "slack",
		Mode:   "single",
		APIURL: "https://slack.com/api/auth.test",
		Request: RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${SECRET}",
				"X-API-Key":     "${INVALID_VAR}",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want invalid variables in single mode error")
	}
}

func TestSecretTemplate_Validate_MultipartWithSecretVariable(t *testing.T) {
	template := SecretTemplate{
		Name:   "stripe",
		Mode:   "multipart",
		APIURL: "https://api.stripe.com/v1/customers",
		RequiredVariables: []string{
			"API_KEY",
		},
		Request: RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${SECRET}",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want SECRET not allowed in multipart mode error")
	}
}

func TestSecretTemplate_Validate_MultipartUndefinedVariables(t *testing.T) {
	template := SecretTemplate{
		Name:   "digitalocean",
		Mode:   "multipart",
		APIURL: "${BASE_URL}/v2/account",
		RequiredVariables: []string{
			"BASE_URL",
		},
		Request: RequestConfig{
			Headers: map[string]string{
				"Authorization": "Bearer ${UNDEFINED_VAR}",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want undefined variables error")
	}
}

func TestSecretTemplate_Validate_MultipartUnusedRequiredVariables(t *testing.T) {
	template := SecretTemplate{
		Name:   "ghost",
		Mode:   "multipart",
		APIURL: "${BASE_URL}/ghost/api/content/posts/",
		RequiredVariables: []string{
			"BASE_URL",
			"UNUSED_VAR",
		},
		Request: RequestConfig{
			Headers: map[string]string{
				"User-Agent": "archer/1.0",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want unused required variables error")
	}
}

func TestSecretTemplate_Validate_RequestConfigError(t *testing.T) {
	data := "some data"
	template := SecretTemplate{
		Name:   "digitalocean",
		Mode:   "single",
		APIURL: "https://api.digitalocean.com/v2/account",
		Request: RequestConfig{
			Data:     &data,
			JSONData: map[string]any{"key": "value"},
			Headers: map[string]string{
				"Authorization": "Bearer ${SECRET}",
			},
		},
	}

	err := template.Validate()

	if err == nil {
		t.Error("Validate() error = nil, want request config validation error")
	}
}

func TestSecretTemplate_Validate_VariableInURL(t *testing.T) {
	template := SecretTemplate{
		Name:   "heroku",
		Mode:   "multipart",
		APIURL: "https://${DOMAIN}/api/${VERSION}",
		RequiredVariables: []string{
			"DOMAIN",
			"VERSION",
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestSecretTemplate_Validate_VariableInQueryParams(t *testing.T) {
	template := SecretTemplate{
		Name:   "ghost",
		Mode:   "multipart",
		APIURL: "https://myblog.ghost.io/ghost/api/content/posts/",
		RequiredVariables: []string{
			"API_KEY",
		},
		Request: RequestConfig{
			QueryParams: map[string]string{
				"key": "${API_KEY}",
			},
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestSecretTemplate_Validate_VariableInData(t *testing.T) {
	data := "grant_type=client_credentials&client_secret=${CLIENT_SECRET}"
	template := SecretTemplate{
		Name:   "twilio",
		Mode:   "multipart",
		APIURL: "https://api.twilio.com/2010-04-01/Accounts",
		RequiredVariables: []string{
			"CLIENT_SECRET",
		},
		Request: RequestConfig{
			Data: &data,
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestSecretTemplate_Validate_VariableInJSONData(t *testing.T) {
	template := SecretTemplate{
		Name:   "mongodb",
		Mode:   "multipart",
		APIURL: "https://cloud.mongodb.com/api/atlas/v1.0/groups",
		RequiredVariables: []string{
			"DATABASE_PASSWORD",
		},
		Request: RequestConfig{
			JSONData: map[string]any{
				"password": "${DATABASE_PASSWORD}",
			},
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestValidationResult_Valid(t *testing.T) {
	result := ValidationResult{
		Valid:   true,
		Message: "Secret is valid",
	}

	if !result.Valid {
		t.Error("Valid = false, want true")
	}

	if result.Error != "" {
		t.Errorf("Error = %q, want empty", result.Error)
	}
}

func TestValidationResult_Invalid(t *testing.T) {
	result := ValidationResult{
		Valid: false,
		Error: "Invalid or expired token",
	}

	if result.Valid {
		t.Error("Valid = true, want false")
	}

	if result.Error != "Invalid or expired token" {
		t.Errorf("Error = %q, want 'Invalid or expired token'", result.Error)
	}
}

func TestSecretTemplate_Validate_AllVariableLocations(t *testing.T) {
	template := SecretTemplate{
		Name:   "datadog",
		Mode:   "multipart",
		APIURL: "${BASE_URL}/api/${VERSION}",
		RequiredVariables: []string{
			"BASE_URL",
			"VERSION",
			"HEADER_VAR",
			"QUERY_VAR",
			"JSON_VAR",
		},
		Request: RequestConfig{
			Headers: map[string]string{
				"X-Custom": "${HEADER_VAR}",
			},
			QueryParams: map[string]string{
				"key": "${QUERY_VAR}",
			},
			JSONData: map[string]any{
				"field": "${JSON_VAR}",
			},
		},
	}

	err := template.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}
