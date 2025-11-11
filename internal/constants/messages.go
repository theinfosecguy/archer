package constants

// Success messages
const (
	SecretValid = "Secret is valid"
)

// Error messages
const (
	TemplateNotFound      = "Template '%s' not found"
	RequestTimeout        = "Request timeout"
	RequestFailed         = "Request failed: %s"
	InvalidJSONResponse   = "Invalid JSON response"
	RequiredFieldNotFound = "Required field '%s' not found"
)

// Template validation messages
const (
	ModeValidationError        = "mode must be either 'single' or 'multipart'"
	MultipartRequiresVariables = "required_variables is mandatory when mode is 'multipart'"
	SingleModeNoVariables      = "required_variables should not be specified when mode is 'single'"
	MutualExclusionError       = "Cannot specify both 'data' and 'json_data'"
	UpperSnakeCaseError        = "Variable '%s' must be in UPPER_SNAKE_CASE format"
	SecretNotAllowedMultipart  = "${SECRET} is not allowed in multipart mode. Use custom variables instead."
	InvalidVariablesSingle     = "In single mode, only ${SECRET} is allowed. Found: %s"
	UndefinedVariables         = "Template uses undefined variables: %s. Add them to required_variables."
	UnusedRequiredVariables    = "Required variables not used in template: %s"
)

// CLI validation messages
const (
	MissingRequiredVariables = "Missing required variables: %s"
	UnexpectedVariables      = "Unexpected variables: %s"
	InvalidVariableFormat    = "Invalid variable format: '%s'. Use --var key=value"
	InvalidKebabCase         = "Variable name '%s' must be in kebab-case format (e.g., 'api-token', 'base-url')"
	VariableNotFound         = "Variable '${%s}' not found in provided variables"
)

// Logging messages
const (
	ValidatorInitialized      = "Validator initialized with templates directory: %s"
	TemplateLoaderInitialized = "Template loader initialized with directory: %s"
	ValidationStarted         = "Starting validation for template '%s' in %s mode"
	ValidationSuccess         = "Validation successful"
	ValidationFailed          = "Validation failed: %s"
	RequestPreparing          = "Preparing %s request to %s"
	RequestSending            = "Sending request with timeout %ds"
	RequestCompleted          = "Request completed with status code %d"
	RequestTimeoutLog         = "Request timed out after %ds"
	RequestFailedLog          = "Request failed: %s"
	TemplateLoaded            = "Template '%s' loaded: %s"
)
