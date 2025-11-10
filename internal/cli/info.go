package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/templates"
	"github.com/theinfosecguy/archer/internal/variables"
)

var infoTemplateFile string

var infoCmd = &cobra.Command{
	Use:   "info TEMPLATE_NAME",
	Short: "Show detailed information about a template",
	Long: `Show detailed information about a template.

Displays comprehensive information about a template's configuration,
required variables, API endpoints, and usage examples.`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	infoCmd.Flags().StringVar(&infoTemplateFile, "template-file", "", "Load template from specific file instead of built-in")
}

func runInfo(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	loader := templates.NewTemplateLoader(constants.DefaultTemplatesDir)
	var templateIdentifier string
	if infoTemplateFile != "" {
		templateIdentifier = infoTemplateFile
	} else {
		templateIdentifier = templateName
	}

	template, err := loader.GetTemplate(templateIdentifier)
	if err != nil {
		return fmt.Errorf("%s Template '%s' not found or invalid", constants.FailureIndicator, templateIdentifier)
	}

	fmt.Printf("Template: %s\n", template.Name)
	fmt.Printf("Description: %s\n", template.Description)
	fmt.Printf("Mode: %s\n", template.Mode)
	fmt.Printf("API URL: %s\n", template.APIURL)
	fmt.Printf("Method: %s\n", template.Method)
	fmt.Println()

	// Show usage information based on mode
	if template.Mode == constants.ModeSingle {
		fmt.Println("Usage:")
		if infoTemplateFile != "" {
			fmt.Printf("  archer validate --template-file %s <secret>\n", infoTemplateFile)
		} else {
			fmt.Printf("  archer validate %s <secret>\n", templateName)
		}
		fmt.Println()
	} else if template.Mode == constants.ModeMultipart {
		fmt.Println("Required Variables:")
		if len(template.RequiredVariables) > 0 {
			for _, varName := range template.RequiredVariables {
				cliName := variables.FormatVarNameForCLI(varName)
				fmt.Printf("  %s (%s %s=<value>)\n", varName, constants.OptVar, cliName)
			}
		}
		fmt.Println()

		fmt.Println("Usage:")
		varExamples := make([]string, 0)
		if len(template.RequiredVariables) > 0 {
			for _, varName := range template.RequiredVariables {
				cliName := variables.FormatVarNameForCLI(varName)
				varExamples = append(varExamples, fmt.Sprintf("%s %s=<value>", constants.OptVar, cliName))
			}
		}

		if infoTemplateFile != "" {
			fmt.Printf("  archer validate --template-file %s %s\n", infoTemplateFile, joinStrings(varExamples, " "))
		} else {
			fmt.Printf("  archer validate %s %s\n", templateName, joinStrings(varExamples, " "))
		}
		fmt.Println()
	}

	fmt.Println("Request Headers:")
	for key, value := range template.Request.Headers {
		maskedValue := variables.MaskVariables(value)
		fmt.Printf("  %s: %s\n", key, maskedValue)
	}

	if len(template.Request.QueryParams) > 0 {
		fmt.Println()
		fmt.Println("Query Parameters:")
		for key, value := range template.Request.QueryParams {
			maskedValue := variables.MaskVariables(value)
			fmt.Printf("  %s: %s\n", key, maskedValue)
		}
	}

	fmt.Println()
	fmt.Printf("Timeout: %ds\n", template.Request.Timeout)
	fmt.Println()

	fmt.Println("Success Criteria:")
	fmt.Printf("  Status Codes: %v\n", template.SuccessCriteria.StatusCode)
	if len(template.SuccessCriteria.RequiredFields) > 0 {
		fmt.Printf("  Required Fields: %s\n", joinStrings(template.SuccessCriteria.RequiredFields, ", "))
	}

	fmt.Println()
	fmt.Println("Error Handling:")
	fmt.Printf("  Max Retries: %d\n", template.ErrorHandling.MaxRetries)
	fmt.Printf("  Retry Delay: %ds\n", template.ErrorHandling.RetryDelay)
	if len(template.ErrorHandling.ErrorMessages) > 0 {
		fmt.Println("  Error Messages:")
		for code, message := range template.ErrorHandling.ErrorMessages {
			fmt.Printf("    %d: %s\n", code, message)
		}
	}

	return nil
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
