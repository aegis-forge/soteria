package statistics

import (
	"gonum.org/v1/gonum/stat"
	"slices"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

func intStatisticsArrayCount(array []models.IntStatistics) models.IntStatistics {
	var totals []int

	for _, intStats := range array {
		totals = append(totals, intStats.Total)
	}

	return BuildIntStatistics(totals)
}

func BuildIntStatistics(array []int) models.IntStatistics {
	floatArray := make([]float64, len(array))

	for _, el := range array {
		floatArray = append(floatArray, float64(el))
	}

	return models.IntStatistics{
		Total:  helpers.ComputeSum(array),
		Min:    slices.Min(array),
		Max:    slices.Max(array),
		Mean:   helpers.ComputeSum(array) / len(array),
		Median: helpers.GetMedian(array),
		Std:    stat.StdDev(floatArray, nil),
	}
}
