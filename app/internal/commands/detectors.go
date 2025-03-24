package commands

import (
	"encoding/json"
	"fmt"
	"tool/app/internal/detector"
	"tool/app/internal/detectors"
)

func Detectors(detects detectors.Detectors) error {
	var detectorsSlice []*detector.Detector

	for _, det := range detects.GetDetectors() {
		detectorsSlice = append(detectorsSlice, det)
	}

	jsonDetect, err := json.Marshal(detectorsSlice)

	if err != nil {
		return err
	}

	fmt.Print(string(jsonDetect))

	return nil
}
