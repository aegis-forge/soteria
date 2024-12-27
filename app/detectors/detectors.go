package detectors

import (
	"errors"
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

func (d *Detectors) EvaluateWorkflow(yamlContent []byte) (map[string][]int, error) {
	var results = make(map[string][]int)

	for key, value := range d.detectorsMap {
		lines, err := value.EvaluateRule(yamlContent)

		if err != nil {
			return nil, err
		}

		if _, ok := results[key]; !ok {
			results[key] = lines
		}
	}

	d.clearDetectors()

	return results, nil
}

func (d *Detectors) clearDetectors() {
	for _, value := range d.detectorsMap {
		value.Rule.ClearResults()
	}
}
