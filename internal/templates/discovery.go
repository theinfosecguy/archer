package templates

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/errors"
)

// DiscoverTemplatesInDirectory discovers all available template files in a directory
func DiscoverTemplatesInDirectory(templatesDir string) ([]string, error) {
	// Check if directory exists
	info, err := os.Stat(templatesDir)
	if os.IsNotExist(err) {
		return nil, &errors.TemplateDirectoryNotFoundError{
			Directory: templatesDir,
		}
	}

	if !info.IsDir() {
		return nil, &errors.TemplateDirectoryNotFoundError{
			Directory: templatesDir,
		}
	}

	templateNames := make(map[string]bool)

	// Walk through directory
	err = filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has valid extension
		ext := filepath.Ext(path)
		if ext == constants.TemplateFileExtension || ext == constants.TemplateFileExtension2 {
			// Get filename without extension
			name := strings.TrimSuffix(filepath.Base(path), ext)
			templateNames[name] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map to sorted slice
	names := make([]string, 0, len(templateNames))
	for name := range templateNames {
		names = append(names, name)
	}
	sort.Strings(names)

	return names, nil
}

// DiscoverTemplates discovers templates from the default templates directory
func DiscoverTemplates(templatesDir string) []string {
	if templatesDir == "" {
		templatesDir = constants.DefaultTemplatesDir
	}

	names, err := DiscoverTemplatesInDirectory(templatesDir)
	if err != nil {
		return []string{}
	}

	return names
}

// GetTemplateIdentifierDisplayName gets a display-friendly name for a template identifier
func GetTemplateIdentifierDisplayName(identifier string) string {
	if IsFilePath(identifier) {
		return strings.TrimSuffix(filepath.Base(identifier), filepath.Ext(identifier))
	}
	return identifier
}
