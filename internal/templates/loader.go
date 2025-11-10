package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/errors"
	"github.com/theinfosecguy/archer/internal/models"
)

// TemplateLoader loads templates from default directory or individual files
type TemplateLoader struct {
	DefaultTemplatesDir string
}

// NewTemplateLoader creates a new template loader
func NewTemplateLoader(defaultTemplatesDir string) *TemplateLoader {
	if defaultTemplatesDir == "" {
		defaultTemplatesDir = constants.DefaultTemplatesDir
	}
	return &TemplateLoader{
		DefaultTemplatesDir: defaultTemplatesDir,
	}
}

// IsFilePath detects if identifier is a file path vs template name
func IsFilePath(identifier string) bool {
	return strings.Contains(identifier, "/") ||
		strings.Contains(identifier, "\\") ||
		strings.HasSuffix(identifier, ".yaml") ||
		strings.HasSuffix(identifier, ".yml")
}

// LoadTemplateFromFile loads a template from a specific file path
func LoadTemplateFromFile(filePath string) (*models.SecretTemplate, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, &errors.TemplateNotFoundError{
			TemplateName: filePath,
		}
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &errors.TemplateLoadError{
			TemplateName: filePath,
			Cause:        err,
		}
	}

	// Parse YAML
	var template models.SecretTemplate
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, &errors.TemplateLoadError{
			TemplateName: filePath,
			Cause:        fmt.Errorf("YAML parsing failed: %w", err),
		}
	}

	// Set defaults
	template.SetDefaults()

	// Validate template
	if err := template.Validate(); err != nil {
		return nil, &errors.TemplateValidationError{
			TemplateName: filePath,
			Message:      err.Error(),
		}
	}

	return &template, nil
}

// LoadTemplateFromDirectory loads a template by name from a templates directory
func LoadTemplateFromDirectory(templateName string, templatesDir string) (*models.SecretTemplate, error) {
	// Check if directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return nil, &errors.TemplateDirectoryNotFoundError{
			Directory: templatesDir,
		}
	}

	// Try .yaml extension first
	templatePath := filepath.Join(templatesDir, templateName+constants.TemplateFileExtension)
	if _, err := os.Stat(templatePath); err == nil {
		return LoadTemplateFromFile(templatePath)
	}

	// Try .yml extension
	templatePath = filepath.Join(templatesDir, templateName+constants.TemplateFileExtension2)
	if _, err := os.Stat(templatePath); err == nil {
		return LoadTemplateFromFile(templatePath)
	}

	return nil, &errors.TemplateNotFoundError{
		TemplateName: templateName,
	}
}

// GetTemplate gets a template by name or file path
func (l *TemplateLoader) GetTemplate(templateIdentifier string) (*models.SecretTemplate, error) {
	if IsFilePath(templateIdentifier) {
		// Direct file path
		return LoadTemplateFromFile(templateIdentifier)
	}

	// Template name - look in default directory
	return LoadTemplateFromDirectory(templateIdentifier, l.DefaultTemplatesDir)
}

// LoadTemplate is an alias for GetTemplate for backward compatibility
func (l *TemplateLoader) LoadTemplate(templateIdentifier string) (*models.SecretTemplate, error) {
	return l.GetTemplate(templateIdentifier)
}
