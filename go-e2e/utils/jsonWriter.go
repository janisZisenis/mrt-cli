package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJSONFile(filePath string, data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
