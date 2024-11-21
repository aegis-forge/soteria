package statistics

import (
	"reflect"
	"regexp"
	"tool/app/internal/models"
)

func environmentArrayCount(environments []models.EnvironmentStatistics) models.EnvironmentStatistics {
	var inherited []int
	var hardcoded []int
	var variables []int
	var counts []int

	for _, environment := range environments {
		inherited = append(inherited, environment.Inherited.Total)
		hardcoded = append(hardcoded, environment.Hardcoded.Total)
		variables = append(variables, environment.Variables.Total)
		counts = append(counts, environment.Inherited.Total+environment.Hardcoded.Total+environment.Variables.Total)
	}

	return models.EnvironmentStatistics{
		Inherited: BuildIntStatistics(inherited),
		Hardcoded: BuildIntStatistics(hardcoded),
		Variables: BuildIntStatistics(variables),
		Count:     BuildIntStatistics(counts),
	}
}

func environmentCount(environment interface{}) models.EnvironmentStatistics {
	regex := regexp.MustCompile(`\$\{\{\s*.+\s*}}`)

	variables := 0
	hardcoded := 0
	inherited := 0

	switch environment := environment.(type) {
	case map[string]interface{}:
		for _, variable := range environment {
			if reflect.TypeOf(variable).Kind().String() == "string" {
				found := regex.FindAllString(variable.(string), -1)

				if len(found) > 0 {
					variables++
				} else {
					hardcoded++
				}
			} else {
				hardcoded++
			}
		}
	case string:
		inherited++
	}

	return models.EnvironmentStatistics{
		Inherited: models.IntStatistics{Total: inherited},
		Hardcoded: models.IntStatistics{Total: hardcoded},
		Variables: models.IntStatistics{Total: variables},
		Count:     models.IntStatistics{Total: variables + hardcoded + inherited},
	}
}
