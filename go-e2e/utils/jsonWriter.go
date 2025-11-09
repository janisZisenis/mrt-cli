package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func TeamConfigWriter(tempDir string, data interface{}) error {
	teamConfig := tempDir + "/team.json"
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if err := os.WriteFile(teamConfig, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
