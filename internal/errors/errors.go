package errors

import "fmt"

// TemplateError represents template-related errors
type TemplateError struct {
	Message string
	Cause   error
}

func (e *TemplateError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// TemplateNotFoundError represents a template not found error
type TemplateNotFoundError struct {
	TemplateName string
}

func (e *TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template not found: %s", e.TemplateName)
}

// TemplateLoadError represents a template loading error
type TemplateLoadError struct {
	TemplateName string
	Cause        error
}

func (e *TemplateLoadError) Error() string {
	return fmt.Sprintf("failed to load template '%s': %v", e.TemplateName, e.Cause)
}

// TemplateValidationError represents a template validation error
type TemplateValidationError struct {
	TemplateName string
	Message      string
}

func (e *TemplateValidationError) Error() string {
	return fmt.Sprintf("template validation failed for '%s': %s", e.TemplateName, e.Message)
}

// TemplateDirectoryNotFoundError represents a templates directory not found error
type TemplateDirectoryNotFoundError struct {
	Directory string
}

func (e *TemplateDirectoryNotFoundError) Error() string {
	return fmt.Sprintf("templates directory not found: %s", e.Directory)
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// CLIError represents CLI-related errors
type CLIError struct {
	Message string
}

func (e *CLIError) Error() string {
	return e.Message
}

// JSONWriteError represents JSON output writing errors
type JSONWriteError struct {
	FilePath string
	Cause    error
}

func (e *JSONWriteError) Error() string {
	return fmt.Sprintf("failed to write JSON to '%s': %v", e.FilePath, e.Cause)
}
