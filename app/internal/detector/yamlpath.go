package detector

import (
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func ResolveYAMLPath(yamlPath string, yamlContent []byte) ([]string, []int, error) {
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
