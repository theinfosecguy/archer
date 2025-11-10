package validator

import (
	"fmt"
	"strings"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/http"
	"github.com/theinfosecguy/archer/internal/logger"
	"github.com/theinfosecguy/archer/internal/models"
	"github.com/theinfosecguy/archer/internal/templates"
	"github.com/theinfosecguy/archer/internal/variables"
)

// SecretValidator validates secrets using configured templates
type SecretValidator struct {
	TemplateLoader *templates.TemplateLoader
	HTTPClient     *http.Client
}

// NewSecretValidator creates a new secret validator
func NewSecretValidator(templatesDir string) *SecretValidator {
	logger.Info("Validator initialized with templates directory: %s", templatesDir)
	return &SecretValidator{
		TemplateLoader: templates.NewTemplateLoader(templatesDir),
		HTTPClient:     http.NewClient(),
	}
}

// ValidateSecret validates a secret using the specified template (single mode)
func (v *SecretValidator) ValidateSecret(templateName string, secret string) (*models.ValidationResult, error) {
	logger.Info("Starting secret validation for template '%s' (mode: single)", templateName)

	template, err := v.TemplateLoader.GetTemplate(templateName)
	if err != nil {
		logger.Info("Validation failed: template '%s' not found in templates directory", templateName)
		return &models.ValidationResult{
			Valid: false,
			Error: fmt.Sprintf(constants.TemplateNotFound, templateName),
		}, nil
	}

	if template.Mode != constants.ModeSingle {
		logger.Info("Template '%s' is not in single mode", templateName)
		return &models.ValidationResult{
			Valid: false,
			Error: fmt.Sprintf("Template '%s' is not in single mode", templateName),
		}, nil
	}

	logger.Debug("Loaded template '%s': %s", template.Name, template.Description)

	// For single mode, create variables map with SECRET
	variablesMap := map[string]string{
		constants.SecretVariableName: secret,
	}

	return v.validateWithTemplate(template, variablesMap)
}

// ValidateSecretMultipart validates secrets using the specified multipart template
func (v *SecretValidator) ValidateSecretMultipart(templateName string, variablesMap map[string]string) (*models.ValidationResult, error) {
	logger.Info("Starting secret validation for template '%s' (mode: multipart)", templateName)

	template, err := v.TemplateLoader.GetTemplate(templateName)
	if err != nil {
		logger.Info("Validation failed: template '%s' not found in templates directory", templateName)
		return &models.ValidationResult{
			Valid: false,
			Error: fmt.Sprintf(constants.TemplateNotFound, templateName),
		}, nil
	}

	if template.Mode != constants.ModeMultipart {
		logger.Info("Template '%s' is not in multipart mode", templateName)
		return &models.ValidationResult{
			Valid: false,
			Error: fmt.Sprintf("Template '%s' is not in multipart mode", templateName),
		}, nil
	}

	logger.Debug("Loaded template '%s': %s", template.Name, template.Description)

	// Validate all required variables are provided
	missingVars := variables.ValidateVariablesProvided(template.RequiredVariables, variablesMap)
	if len(missingVars) > 0 {
		errorMsg := fmt.Sprintf(constants.MissingRequiredVariables, strings.Join(missingVars, ", "))
		logger.Info("Missing required variables: %s", strings.Join(missingVars, ", "))
		return &models.ValidationResult{
			Valid: false,
			Error: errorMsg,
		}, nil
	}

	return v.validateWithTemplate(template, variablesMap)
}

func (v *SecretValidator) validateWithTemplate(template *models.SecretTemplate, vars map[string]string) (*models.ValidationResult, error) {
	// Delegate to HTTP client for request execution
	return v.HTTPClient.ExecuteRequest(template, vars)
}
