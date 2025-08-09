"""Message constants for the Archer secret validation system."""

# Success messages
SECRET_VALID = "Secret is valid"

# Error messages
TEMPLATE_NOT_FOUND = "Template '{template_name}' not found"
REQUEST_TIMEOUT = "Request timeout"
REQUEST_FAILED = "Request failed: {error}"
INVALID_JSON_RESPONSE = "Invalid JSON response"
REQUIRED_FIELD_NOT_FOUND = "Required field '{field_path}' not found"

# Template validation messages
MODE_VALIDATION_ERROR = "mode must be either 'single' or 'multipart'"
MULTIPART_REQUIRES_VARIABLES = "required_variables is mandatory when mode is 'multipart'"
SINGLE_MODE_NO_VARIABLES = "required_variables should not be specified when mode is 'single'"
MUTUAL_EXCLUSION_ERROR = "Cannot specify both 'data' and 'json_data'"
UPPER_SNAKE_CASE_ERROR = "Variable '{var}' must be in UPPER_SNAKE_CASE format"
SECRET_NOT_ALLOWED_MULTIPART = "${SECRET} is not allowed in multipart mode. Use custom variables instead."
INVALID_VARIABLES_SINGLE = "In single mode, only ${{SECRET}} is allowed. Found: {vars}"
UNDEFINED_VARIABLES = "Template uses undefined variables: {vars}. Add them to required_variables."
UNUSED_REQUIRED_VARIABLES = "Required variables not used in template: {vars}"

# CLI validation messages
MISSING_REQUIRED_VARIABLES = "Missing required variables: {vars}"
UNEXPECTED_VARIABLES = "Unexpected variables: {vars}"
INVALID_VARIABLE_FORMAT = "Invalid variable format: '{var}'. Use --var key=value"
INVALID_KEBAB_CASE = "Variable name '{key}' must be in kebab-case format (e.g., 'api-token', 'base-url')"
