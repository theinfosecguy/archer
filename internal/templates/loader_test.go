package templates

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTemplateLoader(t *testing.T) {
	templatesDir := "testdata"
	loader := NewTemplateLoader(templatesDir)

	if loader == nil {
		t.Fatal("NewTemplateLoader() returned nil")
	}

	if loader.DefaultTemplatesDir != templatesDir {
		t.Errorf("DefaultTemplatesDir = %q, want %q", loader.DefaultTemplatesDir, templatesDir)
	}
}

func TestNewTemplateLoader_EmptyDirectory(t *testing.T) {
	loader := NewTemplateLoader("")

	if loader == nil {
		t.Fatal("NewTemplateLoader('') returned nil")
	}

	if loader.DefaultTemplatesDir != "templates" {
		t.Errorf("DefaultTemplatesDir = %q, want 'templates'", loader.DefaultTemplatesDir)
	}
}

func TestIsFilePath_WithForwardSlash(t *testing.T) {
	path := "path/to/template.yaml"

	if !IsFilePath(path) {
		t.Errorf("IsFilePath(%q) = false, want true", path)
	}
}

func TestIsFilePath_WithBackslash(t *testing.T) {
	path := "path\\to\\template.yaml"

	if !IsFilePath(path) {
		t.Errorf("IsFilePath(%q) = false, want true", path)
	}
}

func TestIsFilePath_WithYamlExtension(t *testing.T) {
	path := "github.yaml"

	if !IsFilePath(path) {
		t.Errorf("IsFilePath(%q) = false, want true", path)
	}
}

func TestIsFilePath_WithYmlExtension(t *testing.T) {
	path := "stripe.yml"

	if !IsFilePath(path) {
		t.Errorf("IsFilePath(%q) = false, want true", path)
	}
}

func TestIsFilePath_TemplateName(t *testing.T) {
	name := "github"

	if IsFilePath(name) {
		t.Errorf("IsFilePath(%q) = true, want false", name)
	}
}

func TestIsFilePath_TemplateNameWithHyphen(t *testing.T) {
	name := "new-relic"

	if IsFilePath(name) {
		t.Errorf("IsFilePath(%q) = true, want false", name)
	}
}

func TestLoadTemplateFromFile_ValidSingleMode(t *testing.T) {
	filePath := "testdata/valid_single.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-single" {
		t.Errorf("Name = %q, want 'test-single'", template.Name)
	}

	if template.Mode != "single" {
		t.Errorf("Mode = %q, want 'single'", template.Mode)
	}

	if template.Method != "GET" {
		t.Errorf("Method = %q, want 'GET'", template.Method)
	}

	if template.Request.Timeout != 10 {
		t.Errorf("Timeout = %d, want 10", template.Request.Timeout)
	}
}

func TestLoadTemplateFromFile_ValidMultipartMode(t *testing.T) {
	filePath := "testdata/valid_multipart.yml"

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-multipart" {
		t.Errorf("Name = %q, want 'test-multipart'", template.Name)
	}

	if template.Mode != "multipart" {
		t.Errorf("Mode = %q, want 'multipart'", template.Mode)
	}

	if len(template.RequiredVariables) != 3 {
		t.Errorf("RequiredVariables length = %d, want 3", len(template.RequiredVariables))
	}

	expectedVars := map[string]bool{
		"BASE_URL":      true,
		"API_KEY":       true,
		"CLIENT_SECRET": true,
	}

	for _, v := range template.RequiredVariables {
		if !expectedVars[v] {
			t.Errorf("Unexpected variable: %s", v)
		}
	}
}

func TestLoadTemplateFromFile_FileNotFound(t *testing.T) {
	filePath := "testdata/nonexistent.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err == nil {
		t.Fatal("LoadTemplateFromFile() error = nil, want error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestLoadTemplateFromFile_InvalidYAML(t *testing.T) {
	filePath := "testdata/invalid_yaml.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err == nil {
		t.Fatal("LoadTemplateFromFile() error = nil, want YAML parse error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestLoadTemplateFromFile_InvalidMode(t *testing.T) {
	filePath := "testdata/invalid_mode.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err == nil {
		t.Fatal("LoadTemplateFromFile() error = nil, want validation error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestLoadTemplateFromFile_MultipartMissingVariables(t *testing.T) {
	filePath := "testdata/multipart_no_vars.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err == nil {
		t.Fatal("LoadTemplateFromFile() error = nil, want validation error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestLoadTemplateFromDirectory_WithYamlExtension(t *testing.T) {
	templateName := "valid_single"
	templatesDir := "testdata"

	template, err := LoadTemplateFromDirectory(templateName, templatesDir)

	if err != nil {
		t.Fatalf("LoadTemplateFromDirectory() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-single" {
		t.Errorf("Name = %q, want 'test-single'", template.Name)
	}
}

func TestLoadTemplateFromDirectory_WithYmlExtension(t *testing.T) {
	templateName := "valid_multipart"
	templatesDir := "testdata"

	template, err := LoadTemplateFromDirectory(templateName, templatesDir)

	if err != nil {
		t.Fatalf("LoadTemplateFromDirectory() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-multipart" {
		t.Errorf("Name = %q, want 'test-multipart'", template.Name)
	}
}

func TestLoadTemplateFromDirectory_TemplateNotFound(t *testing.T) {
	templateName := "nonexistent"
	templatesDir := "testdata"

	template, err := LoadTemplateFromDirectory(templateName, templatesDir)

	if err == nil {
		t.Fatal("LoadTemplateFromDirectory() error = nil, want error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestLoadTemplateFromDirectory_DirectoryNotFound(t *testing.T) {
	templateName := "github"
	templatesDir := "nonexistent_directory"

	template, err := LoadTemplateFromDirectory(templateName, templatesDir)

	if err == nil {
		t.Fatal("LoadTemplateFromDirectory() error = nil, want directory not found error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestGetTemplate_WithTemplateName(t *testing.T) {
	loader := NewTemplateLoader("testdata")

	template, err := loader.GetTemplate("valid_single")

	if err != nil {
		t.Fatalf("GetTemplate() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-single" {
		t.Errorf("Name = %q, want 'test-single'", template.Name)
	}
}

func TestGetTemplate_WithFilePath(t *testing.T) {
	loader := NewTemplateLoader("testdata")

	template, err := loader.GetTemplate("testdata/valid_multipart.yml")

	if err != nil {
		t.Fatalf("GetTemplate() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-multipart" {
		t.Errorf("Name = %q, want 'test-multipart'", template.Name)
	}
}

func TestGetTemplate_NotFound(t *testing.T) {
	loader := NewTemplateLoader("testdata")

	template, err := loader.GetTemplate("nonexistent")

	if err == nil {
		t.Fatal("GetTemplate() error = nil, want error")
	}

	if template != nil {
		t.Errorf("template = %v, want nil", template)
	}
}

func TestLoadTemplate_BackwardCompatibility(t *testing.T) {
	loader := NewTemplateLoader("testdata")

	template, err := loader.LoadTemplate("valid_single")

	if err != nil {
		t.Fatalf("LoadTemplate() error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}

	if template.Name != "test-single" {
		t.Errorf("Name = %q, want 'test-single'", template.Name)
	}
}

func TestLoadTemplateFromFile_DefaultsApplied(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "minimal.yaml")

	minimalYAML := `name: sendgrid
description: SendGrid API key validation
api_url: https://api.sendgrid.com/v3/user/profile

request:
  headers:
    Authorization: Bearer ${SECRET}

success_criteria:
  status_code: [200]

error_handling:
  error_messages:
    401: Unauthorized
`

	err := os.WriteFile(filePath, []byte(minimalYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if template.Method != "GET" {
		t.Errorf("Method = %q, want 'GET' (default)", template.Method)
	}

	if template.Mode != "single" {
		t.Errorf("Mode = %q, want 'single' (default)", template.Mode)
	}

	if template.Request.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30 (default)", template.Request.Timeout)
	}

	if template.ErrorHandling.MaxRetries != 0 {
		t.Errorf("MaxRetries = %d, want 0 (default)", template.ErrorHandling.MaxRetries)
	}

	if template.ErrorHandling.RetryDelay != 0 {
		t.Errorf("RetryDelay = %d, want 0 (default)", template.ErrorHandling.RetryDelay)
	}
}

func TestLoadTemplateFromFile_HeadersPresent(t *testing.T) {
	filePath := "testdata/valid_single.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if len(template.Request.Headers) == 0 {
		t.Error("Headers map is empty, want at least one header")
	}

	authHeader, ok := template.Request.Headers["Authorization"]
	if !ok {
		t.Error("Authorization header not found")
	}

	if authHeader != "Bearer ${SECRET}" {
		t.Errorf("Authorization header = %q, want 'Bearer ${SECRET}'", authHeader)
	}
}

func TestLoadTemplateFromFile_SuccessCriteria(t *testing.T) {
	filePath := "testdata/valid_single.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if len(template.SuccessCriteria.StatusCode) == 0 {
		t.Error("StatusCode slice is empty, want at least one status code")
	}

	if template.SuccessCriteria.StatusCode[0] != 200 {
		t.Errorf("StatusCode[0] = %d, want 200", template.SuccessCriteria.StatusCode[0])
	}

	if len(template.SuccessCriteria.RequiredFields) != 2 {
		t.Errorf("RequiredFields length = %d, want 2", len(template.SuccessCriteria.RequiredFields))
	}
}

func TestLoadTemplateFromFile_ErrorHandling(t *testing.T) {
	filePath := "testdata/valid_single.yaml"

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if template.ErrorHandling.MaxRetries != 2 {
		t.Errorf("MaxRetries = %d, want 2", template.ErrorHandling.MaxRetries)
	}

	if template.ErrorHandling.RetryDelay != 1 {
		t.Errorf("RetryDelay = %d, want 1", template.ErrorHandling.RetryDelay)
	}

	if len(template.ErrorHandling.ErrorMessages) == 0 {
		t.Error("ErrorMessages map is empty")
	}

	msg401, ok := template.ErrorHandling.ErrorMessages[401]
	if !ok {
		t.Error("Error message for 401 not found")
	}

	if msg401 != "Invalid token" {
		t.Errorf("Error message 401 = %q, want 'Invalid token'", msg401)
	}
}

func TestLoadTemplateFromFile_JSONData(t *testing.T) {
	filePath := "testdata/valid_multipart.yml"

	template, err := LoadTemplateFromFile(filePath)

	if err != nil {
		t.Fatalf("LoadTemplateFromFile() error = %v, want nil", err)
	}

	if template.Request.JSONData == nil {
		t.Fatal("JSONData is nil, want map")
	}

	if len(template.Request.JSONData) == 0 {
		t.Error("JSONData is empty")
	}

	grantType, ok := template.Request.JSONData["grant_type"]
	if !ok {
		t.Error("grant_type field not found in JSONData")
	}

	if grantType != "client_credentials" {
		t.Errorf("grant_type = %v, want 'client_credentials'", grantType)
	}
}

func TestGetTemplate_RelativePath(t *testing.T) {
	loader := NewTemplateLoader("testdata")

	template, err := loader.GetTemplate("./testdata/valid_single.yaml")

	if err != nil {
		t.Fatalf("GetTemplate() with relative path error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}
}

func TestGetTemplate_AbsolutePath(t *testing.T) {
	loader := NewTemplateLoader("testdata")

	absPath, err := filepath.Abs("testdata/valid_single.yaml")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	template, err := loader.GetTemplate(absPath)

	if err != nil {
		t.Fatalf("GetTemplate() with absolute path error = %v, want nil", err)
	}

	if template == nil {
		t.Fatal("template is nil")
	}
}
