package detectors

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"tool/app/internal/detector"
	"tool/app/internal/detectors/simple"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

var detectorsMap = map[string]map[string]*detector.Detector{
	"simple": {
		"no-hash-version-pin":      &simple.NoHashVersionPin,
		"coarse-permission":        &simple.CoarsePermission,
		"caching-in-release":       &simple.CachingInRelease,
		"bad-local-environment":    &simple.BadLocalEnvironment,
		"bad-github-context":       &simple.BadGithubContext,
		"global-secret":            &simple.GlobalSecret,
		"self-hosted-runner":       &simple.SelfHostedRunner,
		"unsafe-artifact-download": &simple.UnsafeArtifactDownload,
	},
}

type Detectors struct {
	detectorsMap    map[string]*detector.Detector
	detectorResults map[string][]int
}

func (d *Detectors) Init(config models.Config) {
	d.detectorsMap = make(map[string]*detector.Detector)

	if !config.Present {
		for groupName, group := range detectorsMap {
			for name, det := range group {
				d.detectorsMap[groupName+"/"+name] = det
			}
		}
	}

	switch config.Detectors.Method {
	case "include":
		for _, det := range config.Detectors.Names {
			group := strings.Split(det, "/")[0]
			name := strings.Split(det, "/")[1]

			if el, ok := detectorsMap[group][name]; ok {
				d.detectorsMap[group+"/"+name] = el
			}
		}
	case "exclude":
		for groupName, group := range detectorsMap {
			for name, det := range group {
				if slices.Contains(config.Detectors.Names, groupName+"/"+name) {
					continue
				}

				d.detectorsMap[name] = det
			}
		}
	}

	log.Print(d.detectorsMap)
}

func (d *Detectors) GetDetector(name string) (*detector.Detector, error) {
	if res, ok := d.detectorsMap[name]; !ok {
		return nil, errors.New("detector '" + name + "' was not found")
	} else {
		return res, nil
	}
}

func (d *Detectors) AddDetector(detector detector.Detector) error {
	if detector.Name == "" {
		return errors.New("detector name should not be empty")
	}

	if _, ok := d.detectorsMap[detector.Name]; ok {
		return errors.New("detector with same name already exists")
	}

	if detector.Info.Severity < 0 || detector.Info.Severity > 5 {
		return errors.New("detector severity should be between 0 and 5")
	}

	d.detectorsMap[detector.Name] = &detector

	return nil
}

func (d *Detectors) EvaluateWorkflow(workflowName string, yamlContent []byte, verbose bool) (map[string][]int, error) {
	var results = make(map[string][]int)
	var severitiesCount = make(map[int]int)

	if verbose {
		fmt.Print(
			"\033[97;1m",
			strings.Repeat("=", len(workflowName))+"\n",
			workflowName+"\n",
			strings.Repeat("=", len(workflowName))+"\n\n",
			"\033[0m",
		)
	} else {
		fmt.Print("\033[97;1m", workflowName, ":\033[0m ")
	}

	for key, value := range d.detectorsMap {
		lines, err := value.EvaluateRule(yamlContent)

		if err != nil {
			return nil, err
		}

		if _, ok := severitiesCount[value.Info.Severity]; ok {
			severitiesCount[value.Info.Severity] += len(value.Rule.GetLines())
		} else {
			severitiesCount[value.Info.Severity] = len(value.Rule.GetLines())
		}

		if verbose {
			value.PrintResults(yamlContent)
		}

		if _, ok := results[key]; !ok {
			results[key] = lines
		}
	}

	i := 0

	if verbose {
		fmt.Print("Results: ")
	}

	for severity, count := range severitiesCount {
		fmt.Print(helpers.ColorMap[severity]+strings.ToTitle(helpers.SeverityMap[severity]),
			"\u001B[0m "+strconv.Itoa(count),
		)

		i++
		if i != len(severitiesCount) {
			fmt.Print("; ")
		} else {
			if verbose {
				fmt.Print("\n")
			}
		}
	}

	fmt.Println()
	d.clearDetectors()

	return results, nil
}

func (d *Detectors) clearDetectors() {
	for _, value := range d.detectorsMap {
		value.Rule.ClearResults()
	}
}
