package cli

import (
	"github.com/spf13/cobra"

	"github.com/theinfosecguy/archer/internal/constants"
)

const banner = `
>>==------ ARCHER ------==>>
   Secret Validation Tool
`

var rootCmd = &cobra.Command{
	Use:   "archer",
	Short: "Archer - validate secrets against APIs using YAML templates",
	Long: banner + `
Archer is a command-line tool for validating API secrets using YAML templates.

Usage patterns:

1) Single-mode templates (simple API tokens)
   # Using environment variable (recommended - secure)
   export ARCHER_SECRET="ghp_xxxxxxxxxxxxxxxxxxxx"
   archer validate github

   # Using command-line argument (not recommended - insecure)
   archer validate github ghp_xxxxxxxxxxxxxxxxxxxx

2) Multipart templates (multiple parameters)
   # Using environment variables (recommended - secure)
   export ARCHER_VAR_BASE_URL="https://myblog.com"
   export ARCHER_VAR_API_TOKEN="xxxxx"
   archer validate ghost

   # Using --var flags (not recommended - insecure)
   archer validate ghost --var base-url=https://myblog.com --var api-token=xxxxx

3) Custom template files
   archer validate myapi --template-file ./custom-api.yaml
   archer validate custom --template-file ./multipart.yaml

4) Template information
   archer list
   archer info github
   archer info --template-file ./custom.yaml

Security Note:
  Passing secrets as command-line arguments exposes them in shell history,
  process lists, and system logs. Always use environment variables for production.

Get started with: archer list`,
	Version: constants.Version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(infoCmd)
}

// SetArgs sets the command args (useful for testing)
func SetArgs(args []string) {
	rootCmd.SetArgs(args)
}
