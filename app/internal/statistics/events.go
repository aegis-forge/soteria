package statistics

func eventsCount(on interface{}) int {
	switch events := on.(type) {
	case string:
		return 1
	case []interface{}:
		return len(events)
	case map[string]interface{}:
		return len(events)
	default:
		return 0
	}
}
