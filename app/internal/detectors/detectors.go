package detectors

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"tool/app/internal/detector"
	dependency_chain_abuse "tool/app/internal/detectors/dependency-chain-abuse"
	improper_artifact_integrity_validation "tool/app/internal/detectors/improper-artifact-integrity-validation"
	insecure_system_configuration "tool/app/internal/detectors/insecure-system-configuration"
	insufficient_pbac "tool/app/internal/detectors/insufficient-pbac"
	poisoned_pipeline_execution "tool/app/internal/detectors/poinsoned-pipeline-execution"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

var detectorsMap = map[string]map[string]*detector.Detector{
	"dependency-chain-abuse": {
		"no-hash-version-pin": &dependency_chain_abuse.NoHashVersionPin,
	},
	"poisoned-pipeline-execution": {
		"conditional-injection":   &poisoned_pipeline_execution.ConditionalInjection,
		"unconditional-injection": &poisoned_pipeline_execution.UnconditionalInjection,
		"pwn-request":             &poisoned_pipeline_execution.PwnRequest,
	},
	"insufficient-pbac": {
		"coarse-permission": &insufficient_pbac.CoarsePermission,
		"global-secret":     &insufficient_pbac.GlobalSecret,
	},
	"insecure-system-configuration": {
		"self-hosted-runner": &insecure_system_configuration.SelfHostedRunner,
	},
	"improper-artifact-integrity-validation": {
		"caching-in-release":       &improper_artifact_integrity_validation.CachingInRelease,
		"unsafe-artifact-download": &improper_artifact_integrity_validation.UnsafeArtifactDownload,
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

			if els, ok := detectorsMap[group]; ok {
				for key, val := range els {
					d.detectorsMap[group+"/"+key] = val
				}
			}
		}
	case "exclude":
		for groupName, group := range detectorsMap {
			if slices.Contains(config.Detectors.Names, groupName+"/*") {
				continue
			}

			for name, det := range group {
				if slices.Contains(config.Detectors.Names, groupName+"/"+name) {
					continue
				}

				d.detectorsMap[name] = det
			}
		}
	}
}

func (d *Detectors) GetDetector(name string) (*detector.Detector, error) {
	if res, ok := d.detectorsMap[name]; !ok {
		return nil, errors.New("detector '" + name + "' was not found")
	} else {
		return res, nil
	}
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
	}

	for key, value := range d.detectorsMap {
		count := 0
		lines, err := value.EvaluateRule(yamlContent)

		if err != nil {
			return nil, err
		}

		if value.CountAll {
			count = len(value.Rule.GetLines())
		} else if !value.CountAll && len(value.Rule.GetLines()) > 0 {
			count = 1
		}

		if _, ok := severitiesCount[value.Info.Severity]; ok {
			severitiesCount[value.Info.Severity] += count
		} else {
			severitiesCount[value.Info.Severity] = count
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

		for severity, count := range severitiesCount {
			fmt.Print(helpers.ColorMap[severity]+strings.ToTitle(helpers.SeverityMap[severity]),
				"\u001B[0m "+strconv.Itoa(count),
			)

			i++
			if i != len(severitiesCount) {
				fmt.Print("; ")
			} else {
				fmt.Print("\n")
			}
		}

		fmt.Println()
	}

	d.clearDetectors()

	return results, nil
}

func (d *Detectors) clearDetectors() {
	for _, value := range d.detectorsMap {
		value.Rule.ClearResults()
	}
}
