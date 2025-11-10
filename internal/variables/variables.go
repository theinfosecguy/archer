package variables

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/errors"
)

// InjectVariables replaces variable placeholders in content with actual values
func InjectVariables(content string, variables map[string]string) string {
	if content == "" {
		return content
	}

	return constants.VariablePattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := match[2 : len(match)-1] // Remove ${ and }
		if value, ok := variables[varName]; ok {
			return value
		}
		// Return original placeholder if variable not found
		return match
	})
}

// MaskVariables masks all variable values in content for safe logging
func MaskVariables(content string) string {
	if content == "" {
		return content
	}

	return constants.VariablePattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := match[2 : len(match)-1] // Remove ${ and }
		return fmt.Sprintf("%s%s%s", constants.MaskedVariablePrefix, varName, constants.MaskedVariableSuffix)
	})
}

// ProcessHeaders processes headers for both request use and masked logging
func ProcessHeaders(headers map[string]string, variables map[string]string) (map[string]string, map[string]string) {
	requestHeaders := make(map[string]string)
	maskedHeaders := make(map[string]string)

	for key, value := range headers {
		requestHeaders[key] = InjectVariables(value, variables)
		maskedHeaders[key] = MaskVariables(value)
	}

	return requestHeaders, maskedHeaders
}

// ProcessQueryParams processes query parameters for both request use and masked logging
func ProcessQueryParams(queryParams map[string]string, variables map[string]string) (map[string]string, map[string]string) {
	if queryParams == nil {
		return nil, nil
	}

	requestParams := make(map[string]string)
	maskedParams := make(map[string]string)

	for key, value := range queryParams {
		requestParams[key] = InjectVariables(value, variables)
		maskedParams[key] = MaskVariables(value)
	}

	return requestParams, maskedParams
}

// ProcessURL processes URL for both request use and masked logging
func ProcessURL(url string, variables map[string]string) (string, string) {
	requestURL := InjectVariables(url, variables)
	maskedURL := MaskVariables(url)
	return requestURL, maskedURL
}

// ProcessData processes data string for both request use and masked logging
func ProcessData(data *string, variables map[string]string) (*string, *string) {
	if data == nil {
		return nil, nil
	}

	requestData := InjectVariables(*data, variables)
	maskedData := MaskVariables(*data)
	return &requestData, &maskedData
}

// ProcessJSONData processes JSON data for both request use and masked logging
func ProcessJSONData(jsonData map[string]any, variables map[string]string) (map[string]any, map[string]any) {
	if jsonData == nil {
		return nil, nil
	}

	requestData := processJSONRecursive(jsonData, variables, false).(map[string]any)
	maskedData := processJSONRecursive(jsonData, variables, true).(map[string]any)

	return requestData, maskedData
}

func processJSONRecursive(obj any, variables map[string]string, mask bool) any {
	switch v := obj.(type) {
	case string:
		if mask {
			return MaskVariables(v)
		}
		return InjectVariables(v, variables)

	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			result[key] = processJSONRecursive(value, variables, mask)
		}
		return result

	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = processJSONRecursive(item, variables, mask)
		}
		return result

	default:
		// For non-string values (numbers, bools, etc.), return as-is
		return obj
	}
}

// GetVariablesFromTemplateContent extracts all variable names from template content
func GetVariablesFromTemplateContent(content string) []string {
	if content == "" {
		return nil
	}

	matches := constants.VariablePattern.FindAllStringSubmatch(content, -1)
	varSet := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			varSet[match[1]] = true
		}
	}

	vars := make([]string, 0, len(varSet))
	for v := range varSet {
		vars = append(vars, v)
	}

	return vars
}

// ValidateVariablesProvided validates that all required variables are provided
// Returns list of missing variables
func ValidateVariablesProvided(requiredVariables []string, providedVariables map[string]string) []string {
	if len(requiredVariables) == 0 {
		return nil
	}

	missing := make([]string, 0)
	for _, varName := range requiredVariables {
		value, ok := providedVariables[varName]
		if !ok || strings.TrimSpace(value) == "" {
			missing = append(missing, varName)
		}
	}

	return missing
}

// ParseVarArgs parses --var key=value arguments into a dictionary
func ParseVarArgs(varArgs []string) (map[string]string, error) {
	variables := make(map[string]string)

	for _, varArg := range varArgs {
		if !strings.Contains(varArg, constants.VariableSeparator) {
			return nil, &errors.ValidationError{
				Message: fmt.Sprintf(constants.InvalidVariableFormat, varArg),
			}
		}

		// Split only on first '=' to allow '=' in values
		parts := strings.SplitN(varArg, constants.VariableSeparator, 2)
		key := parts[0]
		value := parts[1]

		// Validate key format (kebab-case)
		if !constants.KebabCasePattern.MatchString(key) {
			return nil, &errors.ValidationError{
				Message: fmt.Sprintf(constants.InvalidKebabCase, key),
			}
		}

		// Convert kebab-case to UPPER_SNAKE_CASE
		upperSnakeKey := strings.ToUpper(strings.ReplaceAll(key, constants.KebabToSnakeSeparator, constants.SnakeCaseSeparator))
		variables[upperSnakeKey] = value
	}

	return variables, nil
}

// FormatVarNameForCLI converts UPPER_SNAKE_CASE to kebab-case for CLI display
func FormatVarNameForCLI(upperSnakeName string) string {
	return strings.ToLower(strings.ReplaceAll(upperSnakeName, constants.SnakeCaseSeparator, constants.KebabToSnakeSeparator))
}

// MarshalJSONWithVariables marshals data to JSON string with variable injection
func MarshalJSONWithVariables(data map[string]any, variables map[string]string) (string, error) {
	processed := processJSONRecursive(data, variables, false)
	bytes, err := json.Marshal(processed)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
