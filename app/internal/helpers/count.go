package helpers

import (
	"slices"
	"strings"
	"tool/app/internal/models"
)

func CheckPresence(element interface{}) int {
	switch element.(type) {
	case string:
		if strings.Compare(element.(string), "") != 0 {
			return 1
		}
	case models.Container:
		if element.(models.Container).Image != "" {
			return 1
		}
	case map[string]interface{}:
		if len(element.(map[string]interface{})) != 0 {
			return 1
		}
	default:
		return 0
	}

	return 0
}

func ComputeSum(array []int) int {
	sum := 0

	for i := range array {
		sum += array[i]
	}

	return sum
}

func GetMedian(array []int) int {
	var median int

	slices.Sort(array)

	if len(array)%2 == 0 {
		median = (array[len(array)/2] + array[len(array)/2-1]) / 2
	} else {
		median = array[len(array)/2]
	}

	return median
}
