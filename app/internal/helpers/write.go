package helpers

import (
	"encoding/json"
	"os"
)

func WriteJSONToFile(path string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return err
	}

	return nil
}
