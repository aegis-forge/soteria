package statistics

import "tool/app/internal/models"

func eventsCount(on interface{}) models.IntStatistics {
	var total int

	switch events := on.(type) {
	case string:
		total = 1
	case []interface{}:
		total = len(events)
	case map[string]interface{}:
		total = len(events)
	default:
		total = 0
	}

	return models.IntStatistics{Total: total}
}
