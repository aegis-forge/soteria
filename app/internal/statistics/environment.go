package statistics

import (
	"reflect"
	"regexp"
	"tool/app/internal/models"
)

func environmentCount(environment map[string]interface{}) models.EnvironmentStatistics {
	regex := regexp.MustCompile(`\$\{\{\s*.+\s*}}`)

	variables := 0
	hardcoded := 0

	for _, variable := range environment {
		if reflect.TypeOf(variable).Kind().String() == "string" {
			found := regex.FindAllString(variable.(string), -1)

			if len(found) > 0 {
				variables += 1
			} else {
				hardcoded += 1
			}
		}
	}

	return models.EnvironmentStatistics{
		Hardcoded: models.IntStatistics{Total: hardcoded},
		Variables: models.IntStatistics{Total: variables},
		Count:     models.IntStatistics{Total: variables + hardcoded},
	}
}
