package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/theinfosecguy/archer/internal/constants"
)

// RequestConfig represents request configuration for API calls
type RequestConfig struct {
	Headers     map[string]string `yaml:"headers" json:"headers"`
	Timeout     int               `yaml:"timeout" json:"timeout"`
	Data        *string           `yaml:"data,omitempty" json:"data,omitempty"`
	JSONData    map[string]any    `yaml:"json_data,omitempty" json:"json_data,omitempty"`
	QueryParams map[string]string `yaml:"query_params,omitempty" json:"query_params,omitempty"`
}

// Validate validates the request config
func (r *RequestConfig) Validate() error {
	if r.Data != nil && r.JSONData != nil {
		return fmt.Errorf(constants.MutualExclusionError)
	}
	return nil
}

// SuccessCriteria represents success criteria for validating API responses
type SuccessCriteria struct {
	StatusCode     []int    `yaml:"status_code" json:"status_code"`
	RequiredFields []string `yaml:"required_fields,omitempty" json:"required_fields,omitempty"`
}

// ErrorHandling represents error handling configuration
type ErrorHandling struct {
	MaxRetries    int            `yaml:"max_retries" json:"max_retries"`
	RetryDelay    int            `yaml:"retry_delay" json:"retry_delay"`
	ErrorMessages map[int]string `yaml:"error_messages,omitempty" json:"error_messages,omitempty"`
}

// SecretTemplate represents a template for secret validation
type SecretTemplate struct {
	Name              string          `yaml:"name" json:"name"`
	Description       string          `yaml:"description" json:"description"`
	APIURL            string          `yaml:"api_url" json:"api_url"`
	Method            string          `yaml:"method" json:"method"`
	Mode              string          `yaml:"mode,omitempty" json:"mode,omitempty"`
	RequiredVariables []string        `yaml:"required_variables,omitempty" json:"required_variables,omitempty"`
	Request           RequestConfig   `yaml:"request" json:"request"`
	SuccessCriteria   SuccessCriteria `yaml:"success_criteria" json:"success_criteria"`
	ErrorHandling     ErrorHandling   `yaml:"error_handling" json:"error_handling"`
}

// SetDefaults sets default values for the template
func (t *SecretTemplate) SetDefaults() {
	if t.Method == "" {
		t.Method = constants.MethodGet
	}
	if t.Mode == "" {
		t.Mode = constants.ModeSingle
	}
	if t.Request.Timeout == 0 {
		t.Request.Timeout = constants.DefaultTimeout
	}
	if t.ErrorHandling.MaxRetries == 0 {
		t.ErrorHandling.MaxRetries = constants.DefaultMaxRetries
	}
	if t.ErrorHandling.RetryDelay == 0 {
		t.ErrorHandling.RetryDelay = constants.DefaultRetryDelay
	}
}

// Validate validates the template
func (t *SecretTemplate) Validate() error {
	// Validate mode
	if t.Mode != constants.ModeSingle && t.Mode != constants.ModeMultipart {
		return fmt.Errorf(constants.ModeValidationError)
	}

	// Validate required variables format
	for _, v := range t.RequiredVariables {
		if !constants.UpperSnakeCasePattern.MatchString(v) {
			return fmt.Errorf(constants.UpperSnakeCaseError, v)
		}
	}

	// Validate multipart requirements
	if t.Mode == constants.ModeMultipart && len(t.RequiredVariables) == 0 {
		return fmt.Errorf(constants.MultipartRequiresVariables)
	}

	// Validate single mode requirements
	if t.Mode == constants.ModeSingle && len(t.RequiredVariables) > 0 {
		return fmt.Errorf(constants.SingleModeNoVariables)
	}

	// Validate request config
	if err := t.Request.Validate(); err != nil {
		return err
	}

	// Extract all variables used in template
	usedVariables := make(map[string]bool)

	// Check URL
	for _, match := range constants.VariablePattern.FindAllStringSubmatch(t.APIURL, -1) {
		if len(match) > 1 {
			usedVariables[match[1]] = true
		}
	}

	// Check headers
	for _, value := range t.Request.Headers {
		for _, match := range constants.VariablePattern.FindAllStringSubmatch(value, -1) {
			if len(match) > 1 {
				usedVariables[match[1]] = true
			}
		}
	}

	// Check query params
	for _, value := range t.Request.QueryParams {
		for _, match := range constants.VariablePattern.FindAllStringSubmatch(value, -1) {
			if len(match) > 1 {
				usedVariables[match[1]] = true
			}
		}
	}

	// Check data
	if t.Request.Data != nil {
		for _, match := range constants.VariablePattern.FindAllStringSubmatch(*t.Request.Data, -1) {
			if len(match) > 1 {
				usedVariables[match[1]] = true
			}
		}
	}

	// Check JSON data
	if t.Request.JSONData != nil {
		jsonStr, _ := json.Marshal(t.Request.JSONData)
		for _, match := range constants.VariablePattern.FindAllStringSubmatch(string(jsonStr), -1) {
			if len(match) > 1 {
				usedVariables[match[1]] = true
			}
		}
	}

	// Validate based on mode
	if t.Mode == constants.ModeSingle {
		// Only ${SECRET} should be used
		if len(usedVariables) > 0 {
			delete(usedVariables, constants.SecretVariableName)
			if len(usedVariables) > 0 {
				vars := make([]string, 0, len(usedVariables))
				for v := range usedVariables {
					vars = append(vars, v)
				}
				return fmt.Errorf(constants.InvalidVariablesSingle, strings.Join(vars, ", "))
			}
		}
	} else if t.Mode == constants.ModeMultipart {
		// ${SECRET} should not be used
		if usedVariables[constants.SecretVariableName] {
			return fmt.Errorf(constants.SecretNotAllowedMultipart)
		}

		// All required variables must be used
		if len(t.RequiredVariables) > 0 {
			requiredSet := make(map[string]bool)
			for _, v := range t.RequiredVariables {
				requiredSet[v] = true
			}

			// Check for undefined variables
			undefinedVars := make([]string, 0)
			for v := range usedVariables {
				if !requiredSet[v] {
					undefinedVars = append(undefinedVars, v)
				}
			}
			if len(undefinedVars) > 0 {
				return fmt.Errorf(constants.UndefinedVariables, strings.Join(undefinedVars, ", "))
			}

			// Check for unused required variables
			unusedVars := make([]string, 0)
			for v := range requiredSet {
				if !usedVariables[v] {
					unusedVars = append(unusedVars, v)
				}
			}
			if len(unusedVars) > 0 {
				return fmt.Errorf(constants.UnusedRequiredVariables, strings.Join(unusedVars, ", "))
			}
		}
	}

	return nil
}

// ValidationResult represents the result of a secret validation
type ValidationResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
