package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/logger"
	"github.com/theinfosecguy/archer/internal/models"
	"github.com/theinfosecguy/archer/internal/output"
	"github.com/theinfosecguy/archer/internal/templates"
	"github.com/theinfosecguy/archer/internal/validator"
	"github.com/theinfosecguy/archer/internal/variables"
)

var (
	templateFile string
	varArgs      []string
	verbose      bool
	debug        bool
	outputJSON   string
	jsonOnly     bool
)

var validateCmd = &cobra.Command{
	Use:   "validate TEMPLATE_NAME [SECRET]",
	Short: "Validate a secret using the specified template",
	Long: `Validate a secret using the specified template.

Template mode determines the required arguments and usage pattern.

Single mode:
  # Using environment variable (recommended)
  export ARCHER_SECRET="ghp_xxxxxxxxxxxxxxxxxxxx"
  archer validate github

  # Using command-line argument (shows security warning)
  archer validate github ghp_xxxxxxxxxxxxxxxxxxxx

Multipart mode:
  # Using environment variables (recommended)
  export ARCHER_VAR_BASE_URL="https://myblog.com"
  export ARCHER_VAR_API_TOKEN="xxxxx"
  archer validate ghost

  # Using --var flags (shows security warning)
  archer validate ghost --var base-url=https://myblog.com --var api-token=xxxxx

Security:
  Environment variables prevent secrets from appearing in shell history and process lists.`,
	Args: cobra.MinimumNArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().StringVar(&templateFile, "template-file", "", "Load template from specific file")
	validateCmd.Flags().StringArrayVar(&varArgs, "var", []string{}, "Variable in format key=value (for multipart templates)")
	validateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	validateCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	validateCmd.Flags().StringVarP(&outputJSON, "output-json", "o", "", "Write structured validation result to JSON file")
	validateCmd.Flags().BoolVar(&jsonOnly, "json-only", false, "Suppress normal terminal success output when writing JSON")
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Setup logging based on flags
	if debug {
		logger.SetDebug()
	} else if verbose {
		logger.SetVerbose()
	}

	logger.Info("Starting secret validation process")
	startTime := time.Now().UTC()

	templateName := args[0]
	var secret string
	if len(args) > 1 {
		secret = args[1]
	}

	// Load template to determine mode
	loader := templates.NewTemplateLoader(constants.DefaultTemplatesDir)
	var templateIdentifier string
	if templateFile != "" {
		templateIdentifier = templateFile
	} else {
		templateIdentifier = templateName
	}

	template, err := loader.GetTemplate(templateIdentifier)
	if err != nil {
		errMsg := fmt.Sprintf("Template '%s' not found or invalid", templateIdentifier)
		if outputJSON != "" {
			writeJSONError(outputJSON, templateName, templateFile, nil, nil, startTime, errMsg)
		}
		return fmt.Errorf("%s %s", constants.FailureIndicator, errMsg)
	}

	// Create validator
	v := validator.NewSecretValidator(constants.DefaultTemplatesDir)

	// Validate based on mode
	if template.Mode == constants.ModeSingle {
		return handleSingleMode(v, templateIdentifier, template, secret, startTime)
	} else if template.Mode == constants.ModeMultipart {
		return handleMultipartMode(v, templateIdentifier, template, secret, varArgs, startTime)
	}

	return fmt.Errorf("unsupported template mode: %s", template.Mode)
}

func handleSingleMode(v *validator.SecretValidator, templateIdentifier string, template *models.SecretTemplate, secret string, startTime time.Time) error {
	// Priority 1: Check environment variable
	envSecret := os.Getenv(constants.EnvSecretName)

	var finalSecret string
	var usedCLI bool

	if envSecret != "" {
		finalSecret = envSecret
	} else if secret != "" {
		finalSecret = secret
		usedCLI = true
	} else {
		errMsg := "secret required. Provide via ARCHER_SECRET environment variable or command-line argument"
		if outputJSON != "" {
			vars := map[string]string{constants.SecretVariableName: ""}
			writeJSONError(outputJSON, templateIdentifier, templateFile, template, vars, startTime, errMsg)
		}
		return errors.New(errMsg)
	}

	if len(varArgs) > 0 {
		errMsg := "--var arguments not allowed in single mode"
		if outputJSON != "" {
			vars := map[string]string{constants.SecretVariableName: finalSecret}
			writeJSONError(outputJSON, templateIdentifier, templateFile, template, vars, startTime, errMsg)
		}
		return errors.New(errMsg)
	}

	// Show warning if secret was passed via CLI
	if usedCLI {
		fmt.Fprint(os.Stderr, constants.WarningSecretInCLI)
	}

	// Build variables map for metadata
	vars := map[string]string{constants.SecretVariableName: finalSecret}

	result, err := v.ValidateSecret(templateIdentifier, finalSecret)
	if err != nil {
		if outputJSON != "" {
			writeJSONError(outputJSON, templateIdentifier, templateFile, template, vars, startTime, err.Error())
		}
		return err
	}

	return handleValidationResult(result, template, vars, startTime)
}

func handleMultipartMode(v *validator.SecretValidator, templateIdentifier string, template *models.SecretTemplate, secret string, varArgs []string, startTime time.Time) error {
	if secret != "" {
		errMsg := "secret argument not allowed in multipart mode. Use --var or ARCHER_VAR_* environment variables instead"
		if outputJSON != "" {
			writeJSONError(outputJSON, templateIdentifier, templateFile, template, nil, startTime, errMsg)
		}
		return errors.New(errMsg)
	}

	// Try to get variables from environment first
	envVars := getEnvVariables(template.RequiredVariables)

	var finalVars map[string]string
	var usedCLI bool

	if len(envVars) > 0 {
		// Check if all required variables are present in env
		missingVars := variables.ValidateVariablesProvided(template.RequiredVariables, envVars)
		if len(missingVars) == 0 {
			// All variables available from environment
			finalVars = envVars
		} else if len(varArgs) > 0 {
			// Some missing from env, use CLI args
			parsedVars, err := variables.ParseVarArgs(varArgs)
			if err != nil {
				if outputJSON != "" {
					writeJSONError(outputJSON, templateIdentifier, templateFile, template, envVars, startTime, err.Error())
				}
				return err
			}
			finalVars = parsedVars
			usedCLI = true
		} else {
			errMsg := fmt.Sprintf("missing required variables: %s. Set via ARCHER_VAR_* environment variables or --var flags", strings.Join(missingVars, ", "))
			if outputJSON != "" {
				writeJSONError(outputJSON, templateIdentifier, templateFile, template, envVars, startTime, errMsg)
			}
			return errors.New(errMsg)
		}
	} else if len(varArgs) > 0 {
		// No env vars, use CLI args
		parsedVars, err := variables.ParseVarArgs(varArgs)
		if err != nil {
			if outputJSON != "" {
				writeJSONError(outputJSON, templateIdentifier, templateFile, template, nil, startTime, err.Error())
			}
			return err
		}
		finalVars = parsedVars
		usedCLI = true
	} else {
		errMsg := "--var arguments or ARCHER_VAR_* environment variables required for multipart mode template"
		if outputJSON != "" {
			writeJSONError(outputJSON, templateIdentifier, templateFile, template, nil, startTime, errMsg)
		}
		return errors.New(errMsg)
	}

	// Show warning if variables were passed via CLI
	if usedCLI {
		fmt.Fprint(os.Stderr, constants.WarningSecretInCLI)
	}

	result, err := v.ValidateSecretMultipart(templateIdentifier, finalVars)
	if err != nil {
		if outputJSON != "" {
			writeJSONError(outputJSON, templateIdentifier, templateFile, template, finalVars, startTime, err.Error())
		}
		return err
	}

	return handleValidationResult(result, template, finalVars, startTime)
}

// getEnvVariables retrieves variables from ARCHER_VAR_* environment variables
func getEnvVariables(requiredVariables []string) map[string]string {
	envVars := make(map[string]string)

	for _, varName := range requiredVariables {
		envKey := constants.EnvVarPrefix + varName
		if value := os.Getenv(envKey); value != "" {
			envVars[varName] = value
		}
	}

	return envVars
}

func handleValidationResult(result *models.ValidationResult, template *models.SecretTemplate, vars map[string]string, startTime time.Time) error {
	// Write JSON output if requested
	if outputJSON != "" {
		if err := writeJSONOutput(outputJSON, result, template, vars, startTime); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write JSON output: %v\n", err)
		}
	}

	// Handle terminal output
	if result.Valid {
		// Only show success message if not in json-only mode
		if !(outputJSON != "" && jsonOnly) {
			fmt.Printf("%s %s\n", constants.SuccessIndicator, result.Message)
		}
		return nil
	}

	// Always show errors even in json-only mode
	return fmt.Errorf("%s %s", constants.FailureIndicator, result.Error)
}

// writeJSONOutput writes successful validation result to JSON file
func writeJSONOutput(filepath string, result *models.ValidationResult, template *models.SecretTemplate, vars map[string]string, startTime time.Time) error {
	endTime := time.Now().UTC()

	// Build masked artifacts
	maskedURL, maskedHeaders := buildMaskedArtifacts(template, vars)

	// Get variable names (without values)
	varNames := make([]string, 0, len(vars))
	for k := range vars {
		varNames = append(varNames, k)
	}

	// Determine source
	source := "builtin"
	if templateFile != "" {
		source = "file"
	}

	// Build request metadata
	requestMeta := models.ValidationRequestMeta{
		Template:             templateFile,
		ResolvedTemplateName: &template.Name,
		Mode:                 &template.Mode,
		Source:               &source,
		Method:               &template.Method,
		APIURLMasked:         &maskedURL,
		HeadersMasked:        maskedHeaders,
		QueryParamsMasked:    nil, // TODO: implement if needed
		VariablesProvided:    varNames,
		StartedAt:            startTime,
		FinishedAt:           endTime,
		DurationMS:           float64(endTime.Sub(startTime).Milliseconds()),
	}

	if templateFile == "" {
		requestMeta.Template = template.Name
	}

	// Build response metadata
	responseMeta := models.ValidationResponseMeta{
		StatusCode:            nil, // TODO: capture from HTTP response
		RequiredFieldsChecked: template.SuccessCriteria.RequiredFields,
		FailedRequiredField:   nil,
		Error:                 nil,
	}

	// Build final JSON structure
	jsonOutput := &models.ValidationResultJSON{
		Command: "validate",
		Version: constants.Version,
		Valid:   result.Valid,
		Request: requestMeta,
		Response: responseMeta,
	}

	if result.Valid {
		jsonOutput.Message = &result.Message
	} else {
		jsonOutput.Error = &result.Error
		responseMeta.Error = &result.Error
	}

	return output.WriteJSONFile(filepath, jsonOutput)
}

// writeJSONError writes error validation result to JSON file
func writeJSONError(filepath string, templateName string, templateFilePath string, template *models.SecretTemplate, vars map[string]string, startTime time.Time, errorMsg string) {
	endTime := time.Now().UTC()

	// Get variable names (without values)
	var varNames []string
	if vars != nil {
		varNames = make([]string, 0, len(vars))
		for k := range vars {
			varNames = append(varNames, k)
		}
	}

	// Build request metadata
	var resolvedName, mode, method *string
	var source *string
	var maskedURL *string
	var maskedHeaders map[string]string

	if template != nil {
		resolvedName = &template.Name
		mode = &template.Mode
		method = &template.Method
		src := "builtin"
		if templateFilePath != "" {
			src = "file"
		}
		source = &src

		// Build masked artifacts if we have vars
		if vars != nil {
			url, headers := buildMaskedArtifacts(template, vars)
			maskedURL = &url
			maskedHeaders = headers
		}
	}

	templateID := templateName
	if templateFilePath != "" {
		templateID = templateFilePath
	}

	requestMeta := models.ValidationRequestMeta{
		Template:             templateID,
		ResolvedTemplateName: resolvedName,
		Mode:                 mode,
		Source:               source,
		Method:               method,
		APIURLMasked:         maskedURL,
		HeadersMasked:        maskedHeaders,
		QueryParamsMasked:    nil,
		VariablesProvided:    varNames,
		StartedAt:            startTime,
		FinishedAt:           endTime,
		DurationMS:           float64(endTime.Sub(startTime).Milliseconds()),
	}

	// Build response metadata
	responseMeta := models.ValidationResponseMeta{
		Error: &errorMsg,
	}

	// Build final JSON structure
	jsonOutput := &models.ValidationResultJSON{
		Command:  "validate",
		Version:  constants.Version,
		Valid:    false,
		Error:    &errorMsg,
		Request:  requestMeta,
		Response: responseMeta,
	}

	// Ignore errors writing JSON error output
	_ = output.WriteJSONFile(filepath, jsonOutput)
}

// buildMaskedArtifacts builds masked URL and headers for JSON output
func buildMaskedArtifacts(template *models.SecretTemplate, vars map[string]string) (string, map[string]string) {
	_, maskedURL := variables.ProcessURL(template.APIURL, vars)
	_, maskedHeaders := variables.ProcessHeaders(template.Request.Headers, vars)
	return maskedURL, maskedHeaders
}
