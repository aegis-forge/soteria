package commands

import (
	"encoding/json"
	"fmt"
	"tool/app/internal/detectors"
)

func Detectors(detects detectors.Detectors) error {
	var detectorsSlice []detectors.DetectorOutput

	for _, det := range detects.GetDetectorsWithCategory() {
		detectorsSlice = append(detectorsSlice, det)
	}

	jsonDetect, err := json.Marshal(detectorsSlice)

	if err != nil {
		return err
	}

	fmt.Print(string(jsonDetect))

	return nil
}
