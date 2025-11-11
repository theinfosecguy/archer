package output

import (
	"encoding/json"
	"os"

	"github.com/theinfosecguy/archer/internal/models"
)

// WriteJSONFile writes ValidationResultJSON to a file with pretty formatting
func WriteJSONFile(filepath string, result *models.ValidationResultJSON) error {
	// Marshal with indentation for readability
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restrictive permissions (0644)
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
