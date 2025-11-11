package templates

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverTemplatesInDirectory_ValidDirectory(t *testing.T) {
	templatesDir := "testdata"

	names, err := DiscoverTemplatesInDirectory(templatesDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	if len(names) == 0 {
		t.Error("No templates discovered, want at least one")
	}

	expectedTemplates := map[string]bool{
		"valid_single":      true,
		"valid_multipart":   true,
		"invalid_yaml":      true,
		"invalid_mode":      true,
		"multipart_no_vars": true,
	}

	for _, name := range names {
		if !expectedTemplates[name] {
			t.Errorf("Unexpected template: %s", name)
		}
	}
}

func TestDiscoverTemplatesInDirectory_DirectoryNotFound(t *testing.T) {
	templatesDir := "nonexistent_directory"

	names, err := DiscoverTemplatesInDirectory(templatesDir)

	if err == nil {
		t.Fatal("DiscoverTemplatesInDirectory() error = nil, want directory not found error")
	}

	if names != nil {
		t.Errorf("names = %v, want nil", names)
	}
}

func TestDiscoverTemplatesInDirectory_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	names, err := DiscoverTemplatesInDirectory(tempDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	if len(names) != 0 {
		t.Errorf("Discovered %d templates in empty directory, want 0", len(names))
	}
}

func TestDiscoverTemplatesInDirectory_BothExtensions(t *testing.T) {
	tempDir := t.TempDir()

	yamlContent := `name: github
description: GitHub API token validation
api_url: https://api.github.com/user
method: GET
mode: single

request:
  headers:
    Authorization: Bearer ${SECRET}
  timeout: 10

success_criteria:
  status_code: [200]

error_handling:
  max_retries: 0
  retry_delay: 0
`

	err := os.WriteFile(filepath.Join(tempDir, "github.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .yaml file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "stripe.yml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .yml file: %v", err)
	}

	names, err := DiscoverTemplatesInDirectory(tempDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	if len(names) != 2 {
		t.Errorf("Discovered %d templates, want 2", len(names))
	}

	expectedNames := map[string]bool{
		"github": true,
		"stripe": true,
	}

	for _, name := range names {
		if !expectedNames[name] {
			t.Errorf("Unexpected template name: %s", name)
		}
	}
}

func TestDiscoverTemplatesInDirectory_IgnoresNonTemplateFiles(t *testing.T) {
	tempDir := t.TempDir()

	validYAML := `name: gitlab
description: GitLab API token validation
api_url: https://gitlab.com/api/v4/user
method: GET
mode: single

request:
  headers:
    Authorization: Bearer ${SECRET}
  timeout: 10

success_criteria:
  status_code: [200]

error_handling:
  max_retries: 0
  retry_delay: 0
`

	err := os.WriteFile(filepath.Join(tempDir, "gitlab.yaml"), []byte(validYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create gitlab.yaml: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("readme"), 0644)
	if err != nil {
		t.Fatalf("Failed to create readme.txt: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "config.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create config.json: %v", err)
	}

	names, err := DiscoverTemplatesInDirectory(tempDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	if len(names) != 1 {
		t.Errorf("Discovered %d templates, want 1", len(names))
	}

	if names[0] != "gitlab" {
		t.Errorf("Template name = %q, want 'gitlab'", names[0])
	}
}

func TestDiscoverTemplatesInDirectory_SortedOutput(t *testing.T) {
	tempDir := t.TempDir()

	templateContent := `name: bitbucket
description: Bitbucket API token validation
api_url: https://api.bitbucket.org/2.0/user
method: GET
mode: single

request:
  headers:
    Authorization: Bearer ${SECRET}
  timeout: 10

success_criteria:
  status_code: [200]

error_handling:
  max_retries: 0
  retry_delay: 0
`

	templateNames := []string{"zebra", "apple", "mango", "banana"}
	for _, name := range templateNames {
		err := os.WriteFile(filepath.Join(tempDir, name+".yaml"), []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s.yaml: %v", name, err)
		}
	}

	names, err := DiscoverTemplatesInDirectory(tempDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	expected := []string{"apple", "banana", "mango", "zebra"}

	if len(names) != len(expected) {
		t.Fatalf("Got %d templates, want %d", len(names), len(expected))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("names[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestDiscoverTemplatesInDirectory_FileIsDirectory(t *testing.T) {
	tempDir := t.TempDir()

	templateContent := `name: circleci
description: CircleCI API token validation
api_url: https://circleci.com/api/v2/me
method: GET
mode: single

request:
  headers:
    Authorization: Bearer ${SECRET}
  timeout: 10

success_criteria:
  status_code: [200]

error_handling:
  max_retries: 0
  retry_delay: 0
`

	err := os.WriteFile(filepath.Join(tempDir, "circleci.yaml"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create circleci.yaml: %v", err)
	}

	err = os.Mkdir(filepath.Join(tempDir, "subdirectory"), 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	names, err := DiscoverTemplatesInDirectory(tempDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	if len(names) != 1 {
		t.Errorf("Discovered %d templates, want 1", len(names))
	}

	if names[0] != "circleci" {
		t.Errorf("Template name = %q, want 'circleci'", names[0])
	}
}

func TestDiscoverTemplates_ValidDirectory(t *testing.T) {
	templatesDir := "testdata"

	names := DiscoverTemplates(templatesDir)

	if len(names) == 0 {
		t.Error("No templates discovered, want at least one")
	}
}

func TestDiscoverTemplates_EmptyString(t *testing.T) {
	names := DiscoverTemplates("")

	if names == nil {
		t.Error("DiscoverTemplates('') returned nil, want empty slice")
	}
}

func TestDiscoverTemplates_DirectoryNotFound(t *testing.T) {
	templatesDir := "nonexistent_directory"

	names := DiscoverTemplates(templatesDir)

	if len(names) != 0 {
		t.Errorf("Discovered %d templates from nonexistent directory, want 0", len(names))
	}
}

func TestGetTemplateIdentifierDisplayName_TemplateName(t *testing.T) {
	identifier := "github"

	displayName := GetTemplateIdentifierDisplayName(identifier)

	if displayName != "github" {
		t.Errorf("GetTemplateIdentifierDisplayName(%q) = %q, want 'github'", identifier, displayName)
	}
}

func TestGetTemplateIdentifierDisplayName_FilePathWithYaml(t *testing.T) {
	identifier := "path/to/templates/stripe.yaml"

	displayName := GetTemplateIdentifierDisplayName(identifier)

	if displayName != "stripe" {
		t.Errorf("GetTemplateIdentifierDisplayName(%q) = %q, want 'stripe'", identifier, displayName)
	}
}

func TestGetTemplateIdentifierDisplayName_FilePathWithYml(t *testing.T) {
	identifier := "templates/slack.yml"

	displayName := GetTemplateIdentifierDisplayName(identifier)

	if displayName != "slack" {
		t.Errorf("GetTemplateIdentifierDisplayName(%q) = %q, want 'slack'", identifier, displayName)
	}
}

func TestGetTemplateIdentifierDisplayName_AbsolutePath(t *testing.T) {
	identifier := "/home/user/templates/openai.yaml"

	displayName := GetTemplateIdentifierDisplayName(identifier)

	if displayName != "openai" {
		t.Errorf("GetTemplateIdentifierDisplayName(%q) = %q, want 'openai'", identifier, displayName)
	}
}

func TestGetTemplateIdentifierDisplayName_WindowsPath(t *testing.T) {
	identifier := "templates/digitalocean.yaml"

	displayName := GetTemplateIdentifierDisplayName(identifier)

	if displayName != "digitalocean" {
		t.Errorf("GetTemplateIdentifierDisplayName(%q) = %q, want 'digitalocean'", identifier, displayName)
	}
}

func TestGetTemplateIdentifierDisplayName_RelativePath(t *testing.T) {
	identifier := "./custom/heroku.yml"

	displayName := GetTemplateIdentifierDisplayName(identifier)

	if displayName != "heroku" {
		t.Errorf("GetTemplateIdentifierDisplayName(%q) = %q, want 'heroku'", identifier, displayName)
	}
}

func TestDiscoverTemplatesInDirectory_DuplicateNamesDifferentExtensions(t *testing.T) {
	tempDir := t.TempDir()

	templateContent := `name: travis
description: Travis CI API token validation
api_url: https://api.travis-ci.com/user
method: GET
mode: single

request:
  headers:
    Authorization: Bearer ${SECRET}
  timeout: 10

success_criteria:
  status_code: [200]

error_handling:
  max_retries: 0
  retry_delay: 0
`

	err := os.WriteFile(filepath.Join(tempDir, "gitlab.yaml"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create gitlab.yaml: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "gitlab.yml"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create gitlab.yml: %v", err)
	}

	names, err := DiscoverTemplatesInDirectory(tempDir)

	if err != nil {
		t.Fatalf("DiscoverTemplatesInDirectory() error = %v, want nil", err)
	}

	if len(names) != 1 {
		t.Errorf("Discovered %d templates, want 1 (deduplicated)", len(names))
	}

	if names[0] != "gitlab" {
		t.Errorf("Template name = %q, want 'gitlab'", names[0])
	}
}
