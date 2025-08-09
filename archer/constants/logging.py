"""Logging-related constants for the Archer secret validation system."""

# Log formats
DEFAULT_LOG_FORMAT = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
TIME_FORMAT = "%H:%M:%S"

# Log message templates
VALIDATOR_INITIALIZED = "SecretValidator initialized with templates directory '{dir}'"
TEMPLATE_LOADER_INITIALIZED = "TemplateLoader initialized with directory '{dir}'"
VALIDATION_STARTED = "Starting secret validation using template '{template}' in {mode} mode"
VALIDATION_SUCCESS = "Secret validation completed successfully - all criteria met"
VALIDATION_FAILED = "Secret validation process failed: {error}"
REQUEST_PREPARING = "Preparing {method} request to {url} with headers: {headers}"
REQUEST_SENDING = "Sending HTTP request with {timeout}s timeout"
REQUEST_COMPLETED = "API request completed with status code {status}"
REQUEST_TIMEOUT_LOG = "API request timed out after {timeout} seconds"
REQUEST_FAILED_LOG = "API request failed with exception: {error}"
TEMPLATE_LOADED = "Template '{name}' loaded successfully: {description}"
VARIABLE_NOT_FOUND = "Variable '{var}' not found in provided variables"
