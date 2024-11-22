package statistics

import "tool/app/internal/models"

func defaultsCount(defaults map[string]interface{}) models.IntStatistics {
	count := 0

	for _, value := range defaults {
		switch value.(type) {
		case map[string]interface{}:
			count += defaultsCount(value.(map[string]interface{})).Total
		default:
			count += 1
		}
	}

	return models.IntStatistics{Total: count}
}
