package helpers

import (
	"os"
	"tool/app/internal/models"
)

func ReadWorkflow(path string) (models.Workflow, error) {
	yamlData, err := os.ReadFile(path)

	if err != nil {
		return models.Workflow{}, err
	}

	var workflow models.Workflow
	if err := UnmarshalWorkflow(yamlData, &workflow); err != nil {
		return models.Workflow{}, err
	}

	return workflow, nil
}
