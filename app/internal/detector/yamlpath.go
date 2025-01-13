package detector

import (
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	"slices"
)

func Resolve(yamlPath string, yamlContent []byte) ([]string, []int, error) {
	path, err := yamlpath.NewPath(yamlPath)

	if err != nil {
		return nil, nil, err
	}

	var node yaml.Node

	if err = yaml.Unmarshal(yamlContent, &node); err != nil {
		return nil, nil, err
	}

	out, err := path.Find(&node)

	if err != nil {
		return nil, nil, err
	}

	var nodeValues []string
	var lines []int

	for _, el := range out {
		nodeValues = append(nodeValues, el.Value)
		lines = append(lines, el.Line)
	}

	return nodeValues, lines, nil
}

func CheckExistence(yamlPath string, toCheck string, yamlContent []byte) (bool, []int, error) {
	res, lines, err := Resolve(yamlPath, yamlContent)

	if err != nil {
		return false, nil, err
	}

	if !slices.Contains(res, toCheck) {
		return false, nil, nil
	}

	return true, []int{lines[slices.Index(res, toCheck)]}, nil
}
