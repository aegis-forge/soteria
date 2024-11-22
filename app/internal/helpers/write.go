package helpers

import (
	"encoding/json"
	"os"
	"tool/app/internal/models"
)

func WriteJSONToFile(path string, data models.Statistics) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return err
	}

	return nil
}
