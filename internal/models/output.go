package models

import "time"

// ValidationRequestMeta represents metadata about the validation request
type ValidationRequestMeta struct {
	Template             string            `json:"template"`                      // Template identifier provided by user (name or file path)
	ResolvedTemplateName *string           `json:"resolved_template_name"`        // Template name as defined inside the template file
	Mode                 *string           `json:"mode"`                          // Template mode if template resolved
	Source               *string           `json:"source"`                        // Where the template was loaded from ("builtin" or "file")
	Method               *string           `json:"method"`                        // HTTP method used for validation request
	APIURLMasked         *string           `json:"api_url_masked"`                // Masked API URL with variables hidden
	HeadersMasked        map[string]string `json:"headers_masked,omitempty"`      // Masked request headers
	QueryParamsMasked    map[string]string `json:"query_params_masked,omitempty"` // Masked query parameters
	VariablesProvided    []string          `json:"variables_provided,omitempty"`  // Names of variables provided (values omitted)
	StartedAt            time.Time         `json:"started_at"`                    // UTC timestamp when validation started
	FinishedAt           time.Time         `json:"finished_at"`                   // UTC timestamp when validation finished
	DurationMS           float64           `json:"duration_ms"`                   // Total duration in milliseconds
}

// ValidationResponseMeta represents metadata about the validation response
type ValidationResponseMeta struct {
	StatusCode            *int     `json:"status_code,omitempty"`             // HTTP status code returned by the endpoint if request executed
	RequiredFieldsChecked []string `json:"required_fields_checked,omitempty"` // List of JSONPath fields checked, if any
	FailedRequiredField   *string  `json:"failed_required_field,omitempty"`   // First required field that was missing, if applicable
	Error                 *string  `json:"error,omitempty"`                   // Low-level error encountered before or during request execution
}

// ValidationResultJSON represents the top-level JSON output for validate command
type ValidationResultJSON struct {
	Command  string                 `json:"command"`           // Always "validate"
	Version  string                 `json:"version"`           // Archer version
	Valid    bool                   `json:"valid"`             // Indicates whether the secret validation succeeded
	Message  *string                `json:"message,omitempty"` // Success message when valid is true
	Error    *string                `json:"error,omitempty"`   // Error message when valid is false
	Request  ValidationRequestMeta  `json:"request"`           // Request metadata
	Response ValidationResponseMeta `json:"response"`          // Response metadata
}
