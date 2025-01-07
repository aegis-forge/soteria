package detectors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tool/app/internal/detector"
)

type Detectors struct {
	detectorsMap    map[string]*detector.Detector
	detectorResults map[string][]int
}

func (d *Detectors) Init() {
	d.detectorsMap = map[string]*detector.Detector{
		"no-hash-version-pin":      &NoHashVersionPin,
		"coarse-permissions":       &CoarsePermissions,
		"caching-in-release":       &CachingInRelease,
		"bad-local-environment":    &BadLocalEnvironment,
		"bad-github-context":       &BadGithubContext,
		"global-secrets":           &GlobalSecret,
		"self-hosted-runner":       &SelfHostedRunner,
		"unsafe-artifact-download": &UnsafeArtifactDownload,
	}
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
		fmt.Print(detector.ColorMap[severity]+strings.ToTitle(detector.SeverityMap[severity]),
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
