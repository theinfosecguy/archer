package constants

import "regexp"

// Variable handling
const (
	SecretVariableName = "SECRET"
	VariablePrefix     = "${"
	VariableSuffix     = "}"
)

// Variable patterns
var (
	VariablePattern       = regexp.MustCompile(`\$\{([^}]+)\}`)
	UpperSnakeCasePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	KebabCasePattern      = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
)

// Variable formatting
const (
	VariableSeparator      = "="
	VariableSeparatorCount = 1
	KebabToSnakeSeparator  = "-"
	SnakeCaseSeparator     = "_"
)

// Security
const (
	MaskedVariablePrefix = "***"
	MaskedVariableSuffix = "***"
)
