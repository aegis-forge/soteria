package statistics

import (
	"strconv"
	"strings"
	"tool/app/internal/detectors"
	"tool/app/internal/helpers"
)

func computePermissions(yamlContent []byte, workflows map[string]Group) error {
	if exists, _, err := checkIfExists("$[*]~", "permissions", yamlContent); err == nil {
		groupOut := Group{}
		groupOut.AddManually([]string{}, 0)

		if exists {
			if full, _, err := checkIfExists("$.permissions", "*", yamlContent); err == nil {
				group := Group{}

				if !full {
					group.AddManually([]string{}, 1)
				} else {
					group.AddManually([]string{}, 0)
				}

				workflows["permissions.coarse.none"] = group
			}
		} else {
			groupOut.AddManually([]string{}, 1)
		}

		workflows["permissions.coarse.default"] = groupOut
	} else {
		return err
	}

	return nil
}

func computeWorkflows(yamlContent []byte, workflowName string) (map[string]Group, error) {
	workflow := map[string]Group{}
	toComputeYamlPath := map[string][]string{
		"events":                   {"$.on[*]~", "$.on[*]"},
		"permissions":              {"$.permissions[*]", "$.permissions"},
		"permissions.coarse.read":  {"$.permissions[?(@=='read-all')]"},
		"permissions.coarse.write": {"$.permissions[?(@=='write-all')]"},
		"permissions.fine.read":    {"$.permissions..[?(@=='read')]"},
		"permissions.fine.write":   {"$.permissions..*[?(@=='write')]"},
		"permissions.fine.none":    {"$.permissions..*[?(@=='none')]"},
		"environment":              {"$.env[*]", "$.jobs..env"},
		"environment.inherited":    {"$.env[?(@=='inherited')]"},
		"environment.variables":    {`$.env..*[?(@=~/\$\{\{\s*.+\s*}}/)]`},
		"defaults":                 {"$.defaults..[*]~"},
		"jobs":                     {"$.jobs[*]~"},
	}

	group := Group{Occurrences: []string{workflowName}, Frequencies: 1}
	workflow["count"] = group

	for name, yamlPaths := range toComputeYamlPath {
		group = Group{}

		if err := group.AddOccurrences(yamlPaths, yamlContent); err != nil {
			return nil, err
		}

		workflow[name] = group
	}

	err := computePermissions(yamlContent, workflow)

	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func computeJobs(yamlContent []byte) (map[string]Group, error) {
	jobs := map[string]Group{}
	toComputeYamlPath := map[string][]string{
		"permissions":              {"$.jobs..permissions[*]", "$.jobs..permissions"},
		"permissions.coarse.read":  {"$.jobs..permissions[?(@=='read-all')]"},
		"permissions.coarse.write": {"$.jobs..permissions[?(@=='write-all')]"},
		"permissions.fine.read":    {"$.jobs..permissions..[?(@=='read')]"},
		"permissions.fine.write":   {"$.jobs..permissions..*[?(@=='write')]"},
		"permissions.fine.none":    {"$.jobs..permissions..*[?(@=='none')]"},
		"environment":              {"$.jobs..env[*]", "$.jobs..env"},
		"environment.inherited":    {"$.jobs..env[?(@=='inherited')]"},
		"environment.variables":    {`$.jobs..env..*[?(@=~/\$\{\{\s*.+\s*}}/)]`},
		"needs":                    {"$.jobs..needs[*]", "$.jobs..needs"},
		"if":                       {"$.jobs..if"},
		"runs-on":                  {"$.jobs..runs-on[*]", "$.jobs..runs-on"},
		"local-env":                {"$.jobs..environment[*]", "$.jobs..environment"},
		"defaults":                 {"$.jobs..defaults..[*]~"},
		"services":                 {"$.jobs..services..[*]~"},
		"uses":                     {"$.jobs..uses[?(@=~/workflows/)]"},
		"secrets":                  {"$.jobs..secrets[*]", "$.jobs..secrets"},
	}

	for name, yamlPaths := range toComputeYamlPath {
		group := Group{}

		if err := group.AddOccurrences(yamlPaths, yamlContent); err != nil {
			return nil, err
		}

		jobs[name] = group
	}

	if found, times, err := checkIfExists("$.jobs..[*]~", "container", yamlContent); err == nil {
		group := Group{}

		if found {
			group.AddManually([]string{}, times)
		}

		jobs["containers"] = group
	} else {
		return nil, err
	}

	jobs["steps"] = Group{
		Occurrences: []string{},
		Frequencies: CountOccurrences("$.jobs..steps[*]", yamlContent),
	}

	err := computePermissions(yamlContent, jobs)

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func computeDetectors(yamlContent []byte, lines map[string][]int, path string, detects detectors.Detectors) (map[string]Group, map[string]Group, error) {
	severities := map[string][]string{}

	severitiesConv := map[string]Group{}
	frequencies := map[string]Group{}

	splitPath := strings.Split(path, "/")
	cleanedPath := strings.Join(splitPath[len(splitPath)-2:], "/")

	for detector, occurrences := range lines {
		if detectorObject, err := detects.GetDetector(detector); err == nil {
			if len(occurrences) > 0 {
				var occs []string

				severity := detectorObject.GetSeverity()
				groupFreq := Group{Workflow: cleanedPath}
				freq := 0

				for _, occurrence := range occurrences {
					line, err := helpers.ReadLine(strings.NewReader(string(yamlContent)), occurrence)

					if err != nil {
						return nil, nil, err
					}

					line = strconv.Itoa(occurrence) + " " + strings.TrimSpace(line)
					occs = append(occs, line)
					freq++
				}

				if detectorObject.CountAll {
					for _, line := range occs {
						if _, ok := severities[severity]; !ok {
							severities[severity] = []string{line}
						} else {
							severities[severity] = append(severities[severity], line)
						}
					}
				} else {
					if _, ok := severities[severity]; !ok {
						severities[severity] = []string{occs[0]}
					} else {
						severities[severity] = append(severities[severity], occs[0])
					}
				}

				if detectorObject.CountAll {
					groupFreq.AddManually(occs, freq)
					frequencies[detectorObject.Name] = groupFreq
				} else {
					groupFreq.AddManually(occs, 1)
					frequencies[detectorObject.Name] = groupFreq
				}
			}
		} else {
			return nil, nil, err
		}
	}

	for key, value := range severities {
		severitiesConv[strings.ToLower(key)] = Group{
			Occurrences: value,
			Frequencies: len(value),
		}
	}

	return severitiesConv, frequencies, nil
}
