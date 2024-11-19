package helpers

import (
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"tool/app/internal/models"
)

func SetDefaults(workflow *models.Workflow) {
	for key, value := range workflow.Jobs {
		var job = value

		if job.TimeoutMinutes == 0 {
			job.TimeoutMinutes = Constants.TimeoutMinutes
		}

		workflow.Jobs[key] = job
	}
}

func Unmarshal(data []byte, workflow *models.Workflow) error {
	var raw interface{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &workflow})
	if err := decoder.Decode(raw); err != nil {
		return err
	}

	SetDefaults(workflow)

	return nil
}
