package statistics

func defaultsCount(defaults map[string]interface{}) int {
	count := 0

	for _, value := range defaults {
		switch value.(type) {
		case map[string]interface{}:
			count += defaultsCount(value.(map[string]interface{}))
		default:
			count += 1
		}
	}

	return count
}
