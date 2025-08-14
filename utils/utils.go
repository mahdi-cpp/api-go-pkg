package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteData[T any](data T, filepath string) error {

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tempFile := filepath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, filepath)
}

// ToStringJson converts any struct type to a JSON string
func ToStringJson[T any](data T) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}
	return string(jsonBytes), nil
}
