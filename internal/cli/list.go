package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/templates"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Long: `List all available templates from the default templates directory.

Displays all built-in templates with their mode indicators and descriptions.
Use this command to discover available validation templates before validation.`,
	RunE: runList,
}

func runList(cmd *cobra.Command, args []string) error {
	templateNames := templates.DiscoverTemplates(constants.DefaultTemplatesDir)

	if len(templateNames) == 0 {
		fmt.Println("No templates found.")
		return nil
	}

	fmt.Printf("Available templates (%d):\n\n", len(templateNames))

	loader := templates.NewTemplateLoader(constants.DefaultTemplatesDir)

	sort.Strings(templateNames)
	for _, name := range templateNames {
		template, err := loader.GetTemplate(name)
		if err != nil {
			fmt.Printf("  %-15s [%-10s] - %s\n", name, "invalid", "[Invalid template]")
			continue
		}

		mode := template.Mode
		if mode == "" {
			mode = constants.ModeSingle
		}

		fmt.Printf("  %-15s [%-10s] - %s\n", template.Name, mode, template.Description)
	}

	return nil
}
