package variables

import (
	"testing"
)

func TestInjectVariables_SingleVariable(t *testing.T) {
	content := "Bearer ${API_TOKEN}"
	variables := map[string]string{
		"API_TOKEN": "ghp_1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r",
	}

	result := InjectVariables(content, variables)
	expected := "Bearer ghp_1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r"

	if result != expected {
		t.Errorf("InjectVariables() = %q, want %q", result, expected)
	}
}

func TestInjectVariables_MultipleVariables(t *testing.T) {
	content := "Authorization: ${AUTH_TYPE} ${API_KEY}, User-Agent: ${USER_AGENT}"
	variables := map[string]string{
		"AUTH_TYPE":  "Bearer",
		"API_KEY":    "sk_live_51abcdefghijklmnop",
		"USER_AGENT": "archer/1.0",
	}

	result := InjectVariables(content, variables)
	expected := "Authorization: Bearer sk_live_51abcdefghijklmnop, User-Agent: archer/1.0"

	if result != expected {
		t.Errorf("InjectVariables() = %q, want %q", result, expected)
	}
}

func TestInjectVariables_EmptyContent(t *testing.T) {
	content := ""
	variables := map[string]string{
		"API_TOKEN": "ghp_7xV3uT9sR5qP2nM6lK4jI8hG1fE0dC",
	}

	result := InjectVariables(content, variables)

	if result != "" {
		t.Errorf("InjectVariables() with empty content = %q, want empty string", result)
	}
}

func TestInjectVariables_NoVariablePlaceholders(t *testing.T) {
	content := "Authorization: Bearer ghp_hardcoded_2aB4cD6eF8gH0iJ2kL4mN6oP8qR0"
	variables := map[string]string{
		"API_TOKEN": "ghp_dynamic_9xY7wV5uT3sR1qP0nM8lK6jI4hG2",
	}

	result := InjectVariables(content, variables)

	if result != content {
		t.Errorf("InjectVariables() with no placeholders = %q, want %q", result, content)
	}
}

func TestInjectVariables_MissingVariable(t *testing.T) {
	content := "Bearer ${GITHUB_TOKEN}"
	variables := map[string]string{
		"API_TOKEN": "ghp_some_token",
	}

	result := InjectVariables(content, variables)
	expected := "Bearer ${GITHUB_TOKEN}" // Placeholder should remain

	if result != expected {
		t.Errorf("InjectVariables() with missing variable = %q, want %q", result, expected)
	}
}

func TestInjectVariables_NestedBracesAndNumbers(t *testing.T) {
	content := "Token: ${API_KEY_V2_123}"
	variables := map[string]string{
		"API_KEY_V2_123": "test_xoxb_fake_1234567890_test_abcdefghijklmnopqrstuvwx",
	}

	result := InjectVariables(content, variables)
	expected := "Token: test_xoxb_fake_1234567890_test_abcdefghijklmnopqrstuvwx"

	if result != expected {
		t.Errorf("InjectVariables() with complex var name = %q, want %q", result, expected)
	}
}

func TestInjectVariables_SpecialCharactersInValue(t *testing.T) {
	content := "URL: ${BASE_URL}/api/v1"
	variables := map[string]string{
		"BASE_URL": "https://api.stripe.com:443/path?query=value&key=abc",
	}

	result := InjectVariables(content, variables)
	expected := "URL: https://api.stripe.com:443/path?query=value&key=abc/api/v1"

	if result != expected {
		t.Errorf("InjectVariables() with special chars = %q, want %q", result, expected)
	}
}

func TestInjectVariables_EmptyVariableValue(t *testing.T) {
	content := "Prefix ${EMPTY_VAR} Suffix"
	variables := map[string]string{
		"EMPTY_VAR": "",
	}

	result := InjectVariables(content, variables)
	expected := "Prefix  Suffix"

	if result != expected {
		t.Errorf("InjectVariables() with empty value = %q, want %q", result, expected)
	}
}

func TestInjectVariables_URLWithMultipleVariables(t *testing.T) {
	content := "${BASE_URL}/ghost/api/${API_VERSION}/posts/?key=${API_KEY}"
	variables := map[string]string{
		"BASE_URL":    "https://myblog.ghost.io",
		"API_VERSION": "v4",
		"API_KEY":     "abc123def456ghi789",
	}

	result := InjectVariables(content, variables)
	expected := "https://myblog.ghost.io/ghost/api/v4/posts/?key=abc123def456ghi789"

	if result != expected {
		t.Errorf("InjectVariables() URL = %q, want %q", result, expected)
	}
}

func TestMaskVariables_SingleVariable(t *testing.T) {
	content := "Bearer ${GITHUB_TOKEN}"

	result := MaskVariables(content)
	expected := "Bearer ***GITHUB_TOKEN***"

	if result != expected {
		t.Errorf("MaskVariables() = %q, want %q", result, expected)
	}
}

func TestMaskVariables_MultipleVariables(t *testing.T) {
	content := "Authorization: ${AUTH_TYPE} ${SECRET_KEY}, X-API-Key: ${API_KEY}"

	result := MaskVariables(content)
	expected := "Authorization: ***AUTH_TYPE*** ***SECRET_KEY***, X-API-Key: ***API_KEY***"

	if result != expected {
		t.Errorf("MaskVariables() = %q, want %q", result, expected)
	}
}

func TestMaskVariables_EmptyContent(t *testing.T) {
	content := ""

	result := MaskVariables(content)

	if result != "" {
		t.Errorf("MaskVariables() with empty content = %q, want empty string", result)
	}
}

func TestMaskVariables_NoVariablePlaceholders(t *testing.T) {
	content := "Authorization: Bearer api-key-without-placeholders"

	result := MaskVariables(content)

	if result != content {
		t.Errorf("MaskVariables() with no placeholders = %q, want %q", result, content)
	}
}

func TestMaskVariables_PreservesFormat(t *testing.T) {
	content := "POST /api/v1/validate\nAuthorization: Bearer ${SECRET}\nContent-Type: application/json"

	result := MaskVariables(content)
	expected := "POST /api/v1/validate\nAuthorization: Bearer ***SECRET***\nContent-Type: application/json"

	if result != expected {
		t.Errorf("MaskVariables() format = %q, want %q", result, expected)
	}
}

func TestMaskVariables_ComplexURL(t *testing.T) {
	content := "https://${DOMAIN}/api/v1/users/${USER_ID}?token=${ACCESS_TOKEN}"

	result := MaskVariables(content)
	expected := "https://***DOMAIN***/api/v1/users/***USER_ID***?token=***ACCESS_TOKEN***"

	if result != expected {
		t.Errorf("MaskVariables() URL = %q, want %q", result, expected)
	}
}

func TestParseVarArgs_ValidSingleArgument(t *testing.T) {
	varArgs := []string{"api-key=sk_live_abcdef123456"}

	result, err := ParseVarArgs(varArgs)
	if err != nil {
		t.Fatalf("ParseVarArgs() error = %v, want nil", err)
	}

	expected := map[string]string{
		"API_KEY": "sk_live_abcdef123456",
	}

	if result["API_KEY"] != expected["API_KEY"] {
		t.Errorf("ParseVarArgs() = %v, want %v", result, expected)
	}
}

func TestParseVarArgs_ValidMultipleArguments(t *testing.T) {
	varArgs := []string{
		"api-key=sk_live_stripe_key",
		"base-url=https://api.stripe.com",
		"webhook-secret=whsec_abc123",
	}

	result, err := ParseVarArgs(varArgs)
	if err != nil {
		t.Fatalf("ParseVarArgs() error = %v, want nil", err)
	}

	if len(result) != 3 {
		t.Errorf("ParseVarArgs() returned %d variables, want 3", len(result))
	}

	if result["API_KEY"] != "sk_live_stripe_key" {
		t.Errorf("API_KEY = %q, want %q", result["API_KEY"], "sk_live_stripe_key")
	}

	if result["BASE_URL"] != "https://api.stripe.com" {
		t.Errorf("BASE_URL = %q, want %q", result["BASE_URL"], "https://api.stripe.com")
	}

	if result["WEBHOOK_SECRET"] != "whsec_abc123" {
		t.Errorf("WEBHOOK_SECRET = %q, want %q", result["WEBHOOK_SECRET"], "whsec_abc123")
	}
}

func TestParseVarArgs_ValueWithEqualsSign(t *testing.T) {
	varArgs := []string{"connection-string=Server=localhost;Database=mydb;User=admin;Password=pass=word"}

	result, err := ParseVarArgs(varArgs)
	if err != nil {
		t.Fatalf("ParseVarArgs() error = %v, want nil", err)
	}

	expected := "Server=localhost;Database=mydb;User=admin;Password=pass=word"
	if result["CONNECTION_STRING"] != expected {
		t.Errorf("CONNECTION_STRING = %q, want %q", result["CONNECTION_STRING"], expected)
	}
}

func TestParseVarArgs_InvalidFormatMissingEquals(t *testing.T) {
	varArgs := []string{"invalid-format-no-equals"}

	_, err := ParseVarArgs(varArgs)
	if err == nil {
		t.Error("ParseVarArgs() with invalid format should return error")
	}
}

func TestParseVarArgs_InvalidKebabCaseUpperCase(t *testing.T) {
	varArgs := []string{"API_KEY=sk_live_key"}

	_, err := ParseVarArgs(varArgs)
	if err == nil {
		t.Error("ParseVarArgs() with UPPER_SNAKE_CASE should return error")
	}
}

func TestParseVarArgs_InvalidKebabCaseMixedCase(t *testing.T) {
	varArgs := []string{"apiKey=sk_live_key"}

	_, err := ParseVarArgs(varArgs)
	if err == nil {
		t.Error("ParseVarArgs() with camelCase should return error")
	}
}

func TestParseVarArgs_EmptyKey(t *testing.T) {
	varArgs := []string{"=some-value"}

	_, err := ParseVarArgs(varArgs)
	if err == nil {
		t.Error("ParseVarArgs() with empty key should return error")
	}
}

func TestParseVarArgs_EmptyValue(t *testing.T) {
	varArgs := []string{"github-token="}

	result, err := ParseVarArgs(varArgs)
	if err != nil {
		t.Fatalf("ParseVarArgs() with empty value error = %v, want nil", err)
	}

	if result["GITHUB_TOKEN"] != "" {
		t.Errorf("GITHUB_TOKEN = %q, want empty string", result["GITHUB_TOKEN"])
	}
}

func TestParseVarArgs_KebabToSnakeCaseConversion(t *testing.T) {
	varArgs := []string{
		"my-api-token=test_sk_fake_4H9mK3nP7qR2sT6vW8xY1zA4",
		"base-url-v2=https://api.datadog.com",
		"oauth-client-secret=test_cs_fake_9aB2cD5eF8gH1iJ4kL7m",
	}

	result, err := ParseVarArgs(varArgs)
	if err != nil {
		t.Fatalf("ParseVarArgs() error = %v, want nil", err)
	}

	testCases := []struct {
		key      string
		expected string
	}{
		{"MY_API_TOKEN", "test_sk_fake_4H9mK3nP7qR2sT6vW8xY1zA4"},
		{"BASE_URL_V2", "https://api.datadog.com"},
		{"OAUTH_CLIENT_SECRET", "test_cs_fake_9aB2cD5eF8gH1iJ4kL7m"},
	}

	for _, tc := range testCases {
		if result[tc.key] != tc.expected {
			t.Errorf("%s = %q, want %q", tc.key, result[tc.key], tc.expected)
		}
	}
}

func TestValidateVariablesProvided_AllProvided(t *testing.T) {
	requiredVariables := []string{"API_KEY", "BASE_URL", "WEBHOOK_SECRET"}
	providedVariables := map[string]string{
		"API_KEY":        "sk_live_key",
		"BASE_URL":       "https://api.stripe.com",
		"WEBHOOK_SECRET": "whsec_abc123",
	}

	missing := ValidateVariablesProvided(requiredVariables, providedVariables)

	if len(missing) != 0 {
		t.Errorf("ValidateVariablesProvided() = %v, want empty slice", missing)
	}
}

func TestValidateVariablesProvided_SomeMissing(t *testing.T) {
	requiredVariables := []string{"GITHUB_TOKEN", "GITLAB_TOKEN", "BITBUCKET_TOKEN"}
	providedVariables := map[string]string{
		"GITHUB_TOKEN": "ghp_2aB8cD4eF6gH0iJ3kL5mN7oP9qR1sT",
	}

	missing := ValidateVariablesProvided(requiredVariables, providedVariables)

	if len(missing) != 2 {
		t.Errorf("ValidateVariablesProvided() returned %d missing, want 2", len(missing))
	}

	expectedMissing := map[string]bool{
		"GITLAB_TOKEN":    true,
		"BITBUCKET_TOKEN": true,
	}

	for _, varName := range missing {
		if !expectedMissing[varName] {
			t.Errorf("Unexpected missing variable: %s", varName)
		}
	}
}

func TestValidateVariablesProvided_AllMissing(t *testing.T) {
	requiredVariables := []string{"API_KEY", "SECRET_KEY"}
	providedVariables := map[string]string{}

	missing := ValidateVariablesProvided(requiredVariables, providedVariables)

	if len(missing) != 2 {
		t.Errorf("ValidateVariablesProvided() returned %d missing, want 2", len(missing))
	}
}

func TestValidateVariablesProvided_EmptyRequired(t *testing.T) {
	requiredVariables := []string{}
	providedVariables := map[string]string{
		"API_KEY": "sk_key",
	}

	missing := ValidateVariablesProvided(requiredVariables, providedVariables)

	if missing != nil {
		t.Errorf("ValidateVariablesProvided() with no required = %v, want nil", missing)
	}
}

func TestValidateVariablesProvided_WhitespaceValue(t *testing.T) {
	requiredVariables := []string{"SLACK_TOKEN"}
	providedVariables := map[string]string{
		"SLACK_TOKEN": "   ",
	}

	missing := ValidateVariablesProvided(requiredVariables, providedVariables)

	if len(missing) != 1 || missing[0] != "SLACK_TOKEN" {
		t.Errorf("ValidateVariablesProvided() with whitespace = %v, want [SLACK_TOKEN]", missing)
	}
}

func TestValidateVariablesProvided_NilProvided(t *testing.T) {
	requiredVariables := []string{"API_KEY"}
	var providedVariables map[string]string = nil

	missing := ValidateVariablesProvided(requiredVariables, providedVariables)

	if len(missing) != 1 {
		t.Errorf("ValidateVariablesProvided() with nil provided = %v, want [API_KEY]", missing)
	}
}

func TestProcessJSONData_SimpleObject(t *testing.T) {
	jsonData := map[string]any{
		"api_key": "${STRIPE_KEY}",
		"secret":  "${STRIPE_SECRET}",
	}
	variables := map[string]string{
		"STRIPE_KEY":    "sk_live_abc123",
		"STRIPE_SECRET": "whsec_def456",
	}

	requestData, maskedData := ProcessJSONData(jsonData, variables)

	if requestData["api_key"] != "sk_live_abc123" {
		t.Errorf("requestData api_key = %v, want sk_live_abc123", requestData["api_key"])
	}

	if maskedData["api_key"] != "***STRIPE_KEY***" {
		t.Errorf("maskedData api_key = %v, want ***STRIPE_KEY***", maskedData["api_key"])
	}
}

func TestProcessJSONData_NestedObject(t *testing.T) {
	jsonData := map[string]any{
		"credentials": map[string]any{
			"username": "admin",
			"password": "${DATABASE_PASSWORD}",
			"connection": map[string]any{
				"host": "${DB_HOST}",
				"port": 5432,
			},
		},
	}
	variables := map[string]string{
		"DATABASE_PASSWORD": "secure_pass_123",
		"DB_HOST":           "localhost",
	}

	requestData, maskedData := ProcessJSONData(jsonData, variables)

	credentials := requestData["credentials"].(map[string]any)
	if credentials["password"] != "secure_pass_123" {
		t.Errorf("nested password = %v, want secure_pass_123", credentials["password"])
	}

	connection := credentials["connection"].(map[string]any)
	if connection["host"] != "localhost" {
		t.Errorf("nested host = %v, want localhost", connection["host"])
	}

	maskedCredentials := maskedData["credentials"].(map[string]any)
	if maskedCredentials["password"] != "***DATABASE_PASSWORD***" {
		t.Errorf("masked password = %v, want ***DATABASE_PASSWORD***", maskedCredentials["password"])
	}
}

func TestProcessJSONData_Array(t *testing.T) {
	jsonData := map[string]any{
		"tokens": []any{"${TOKEN_1}", "${TOKEN_2}", "${TOKEN_3}"},
	}
	variables := map[string]string{
		"TOKEN_1": "ghp_token_alpha",
		"TOKEN_2": "ghp_token_beta",
		"TOKEN_3": "ghp_token_gamma",
	}

	requestData, maskedData := ProcessJSONData(jsonData, variables)

	tokens := requestData["tokens"].([]any)
	if tokens[0] != "ghp_token_alpha" {
		t.Errorf("token[0] = %v, want ghp_token_alpha", tokens[0])
	}

	maskedTokens := maskedData["tokens"].([]any)
	if maskedTokens[1] != "***TOKEN_2***" {
		t.Errorf("maskedToken[1] = %v, want ***TOKEN_2***", maskedTokens[1])
	}
}

func TestProcessJSONData_MixedTypes(t *testing.T) {
	jsonData := map[string]any{
		"api_key":    "${OPENAI_KEY}",
		"model":      "gpt-4",
		"max_tokens": 1000,
		"stream":     true,
		"metadata":   nil,
	}
	variables := map[string]string{
		"OPENAI_KEY": "sk-proj-abc123",
	}

	requestData, _ := ProcessJSONData(jsonData, variables)

	if requestData["api_key"] != "sk-proj-abc123" {
		t.Errorf("api_key = %v, want sk-proj-abc123", requestData["api_key"])
	}

	if requestData["model"] != "gpt-4" {
		t.Errorf("model = %v, want gpt-4", requestData["model"])
	}

	if requestData["max_tokens"] != 1000 {
		t.Errorf("max_tokens = %v, want 1000", requestData["max_tokens"])
	}

	if requestData["stream"] != true {
		t.Errorf("stream = %v, want true", requestData["stream"])
	}
}

func TestProcessJSONData_EmptyObject(t *testing.T) {
	jsonData := map[string]any{}
	variables := map[string]string{
		"API_KEY": "sk_key",
	}

	requestData, maskedData := ProcessJSONData(jsonData, variables)

	if len(requestData) != 0 {
		t.Errorf("requestData length = %d, want 0", len(requestData))
	}

	if len(maskedData) != 0 {
		t.Errorf("maskedData length = %d, want 0", len(maskedData))
	}
}

func TestProcessJSONData_NullValues(t *testing.T) {
	jsonData := map[string]any{
		"optional_field": nil,
		"required_field": "${API_TOKEN}",
	}
	variables := map[string]string{
		"API_TOKEN": "token_value",
	}

	requestData, _ := ProcessJSONData(jsonData, variables)

	if requestData["optional_field"] != nil {
		t.Errorf("optional_field = %v, want nil", requestData["optional_field"])
	}

	if requestData["required_field"] != "token_value" {
		t.Errorf("required_field = %v, want token_value", requestData["required_field"])
	}
}

func TestProcessJSONData_ArrayOfObjects(t *testing.T) {
	jsonData := map[string]any{
		"servers": []any{
			map[string]any{
				"name": "production",
				"url":  "${PROD_URL}",
			},
			map[string]any{
				"name": "staging",
				"url":  "${STAGING_URL}",
			},
		},
	}
	variables := map[string]string{
		"PROD_URL":    "https://api.production.aws.acme.io",
		"STAGING_URL": "https://api.staging.aws.acme.io",
	}

	requestData, maskedData := ProcessJSONData(jsonData, variables)

	servers := requestData["servers"].([]any)
	prodServer := servers[0].(map[string]any)

	if prodServer["url"] != "https://api.production.aws.acme.io" {
		t.Errorf("production url = %v, want https://api.production.aws.acme.io", prodServer["url"])
	}

	maskedServers := maskedData["servers"].([]any)
	maskedProd := maskedServers[0].(map[string]any)

	if maskedProd["url"] != "***PROD_URL***" {
		t.Errorf("masked production url = %v, want ***PROD_URL***", maskedProd["url"])
	}
}

func TestProcessJSONData_NilInput(t *testing.T) {
	var jsonData map[string]any = nil
	variables := map[string]string{
		"API_KEY": "sk_key",
	}

	requestData, maskedData := ProcessJSONData(jsonData, variables)

	if requestData != nil {
		t.Errorf("requestData = %v, want nil", requestData)
	}

	if maskedData != nil {
		t.Errorf("maskedData = %v, want nil", maskedData)
	}
}

func TestFormatVarNameForCLI(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"API_KEY", "api-key"},
		{"BASE_URL", "base-url"},
		{"OAUTH_CLIENT_SECRET", "oauth-client-secret"},
		{"MY_API_TOKEN_V2", "my-api-token-v2"},
		{"SIMPLE", "simple"},
	}

	for _, tc := range testCases {
		result := FormatVarNameForCLI(tc.input)
		if result != tc.expected {
			t.Errorf("FormatVarNameForCLI(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestGetVariablesFromTemplateContent(t *testing.T) {
	content := "Authorization: Bearer ${GITHUB_TOKEN}, X-API-Version: ${API_VERSION}"

	vars := GetVariablesFromTemplateContent(content)

	if len(vars) != 2 {
		t.Errorf("GetVariablesFromTemplateContent() returned %d vars, want 2", len(vars))
	}

	varMap := make(map[string]bool)
	for _, v := range vars {
		varMap[v] = true
	}

	if !varMap["GITHUB_TOKEN"] {
		t.Error("Missing GITHUB_TOKEN in extracted variables")
	}

	if !varMap["API_VERSION"] {
		t.Error("Missing API_VERSION in extracted variables")
	}
}

func TestGetVariablesFromTemplateContent_EmptyContent(t *testing.T) {
	content := ""

	vars := GetVariablesFromTemplateContent(content)

	if vars != nil {
		t.Errorf("GetVariablesFromTemplateContent('') = %v, want nil", vars)
	}
}

func TestGetVariablesFromTemplateContent_NoVariables(t *testing.T) {
	content := "Authorization: Bearer static-token"

	vars := GetVariablesFromTemplateContent(content)

	if len(vars) != 0 {
		t.Errorf("GetVariablesFromTemplateContent() = %v, want empty", vars)
	}
}

func TestGetVariablesFromTemplateContent_DuplicateVariables(t *testing.T) {
	content := "${API_KEY} and ${API_KEY} and ${SECRET}"

	vars := GetVariablesFromTemplateContent(content)

	if len(vars) != 2 {
		t.Errorf("GetVariablesFromTemplateContent() with duplicates returned %d vars, want 2", len(vars))
	}
}

func TestProcessHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer ${SLACK_TOKEN}",
		"Content-Type":  "application/json",
		"X-Custom":      "${CUSTOM_HEADER}",
	}
	variables := map[string]string{
		"SLACK_TOKEN":   "test_fake_xoxb_1234_5678_abcd",
		"CUSTOM_HEADER": "custom-value",
	}

	requestHeaders, maskedHeaders := ProcessHeaders(headers, variables)

	if requestHeaders["Authorization"] != "Bearer test_fake_xoxb_1234_5678_abcd" {
		t.Errorf("Authorization header = %q, want Bearer test_fake_xoxb_1234_5678_abcd", requestHeaders["Authorization"])
	}

	if maskedHeaders["Authorization"] != "Bearer ***SLACK_TOKEN***" {
		t.Errorf("Masked Authorization = %q, want Bearer ***SLACK_TOKEN***", maskedHeaders["Authorization"])
	}

	if requestHeaders["Content-Type"] != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", requestHeaders["Content-Type"])
	}
}

func TestProcessQueryParams(t *testing.T) {
	queryParams := map[string]string{
		"key":    "${API_KEY}",
		"format": "json",
		"limit":  "10",
	}
	variables := map[string]string{
		"API_KEY": "abc123def456",
	}

	requestParams, maskedParams := ProcessQueryParams(queryParams, variables)

	if requestParams["key"] != "abc123def456" {
		t.Errorf("key param = %q, want abc123def456", requestParams["key"])
	}

	if maskedParams["key"] != "***API_KEY***" {
		t.Errorf("masked key = %q, want ***API_KEY***", maskedParams["key"])
	}

	if requestParams["format"] != "json" {
		t.Errorf("format param = %q, want json", requestParams["format"])
	}
}

func TestProcessQueryParams_NilInput(t *testing.T) {
	var queryParams map[string]string = nil
	variables := map[string]string{
		"API_KEY": "key123",
	}

	requestParams, maskedParams := ProcessQueryParams(queryParams, variables)

	if requestParams != nil {
		t.Errorf("requestParams = %v, want nil", requestParams)
	}

	if maskedParams != nil {
		t.Errorf("maskedParams = %v, want nil", maskedParams)
	}
}

func TestProcessURL(t *testing.T) {
	url := "${BASE_URL}/api/v1/users/${USER_ID}"
	variables := map[string]string{
		"BASE_URL": "https://api.github.com",
		"USER_ID":  "12345",
	}

	requestURL, maskedURL := ProcessURL(url, variables)

	expectedRequest := "https://api.github.com/api/v1/users/12345"
	expectedMasked := "***BASE_URL***/api/v1/users/***USER_ID***"

	if requestURL != expectedRequest {
		t.Errorf("requestURL = %q, want %q", requestURL, expectedRequest)
	}

	if maskedURL != expectedMasked {
		t.Errorf("maskedURL = %q, want %q", maskedURL, expectedMasked)
	}
}

func TestProcessData(t *testing.T) {
	data := "grant_type=client_credentials&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}"
	variables := map[string]string{
		"CLIENT_ID":     "app_client_id",
		"CLIENT_SECRET": "app_client_secret",
	}

	requestData, maskedData := ProcessData(&data, variables)

	expectedRequest := "grant_type=client_credentials&client_id=app_client_id&client_secret=app_client_secret"
	expectedMasked := "grant_type=client_credentials&client_id=***CLIENT_ID***&client_secret=***CLIENT_SECRET***"

	if *requestData != expectedRequest {
		t.Errorf("requestData = %q, want %q", *requestData, expectedRequest)
	}

	if *maskedData != expectedMasked {
		t.Errorf("maskedData = %q, want %q", *maskedData, expectedMasked)
	}
}

func TestProcessData_NilInput(t *testing.T) {
	var data *string = nil
	variables := map[string]string{
		"API_KEY": "key123",
	}

	requestData, maskedData := ProcessData(data, variables)

	if requestData != nil {
		t.Errorf("requestData = %v, want nil", requestData)
	}

	if maskedData != nil {
		t.Errorf("maskedData = %v, want nil", maskedData)
	}
}

func TestMarshalJSONWithVariables(t *testing.T) {
	data := map[string]any{
		"username": "admin",
		"password": "${DATABASE_PASSWORD}",
		"port":     5432,
	}
	variables := map[string]string{
		"DATABASE_PASSWORD": "secure_pass_123",
	}

	result, err := MarshalJSONWithVariables(data, variables)
	if err != nil {
		t.Fatalf("MarshalJSONWithVariables() error = %v, want nil", err)
	}

	// Check that password was injected
	if !contains(result, "secure_pass_123") {
		t.Errorf("MarshalJSONWithVariables() result doesn't contain injected password")
	}

	// Check that it's valid JSON
	if !contains(result, `"username":"admin"`) && !contains(result, `"username": "admin"`) {
		t.Errorf("MarshalJSONWithVariables() result doesn't contain username field")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
