package detector

import (
	"reflect"
	"slices"
	"strings"
)

var severityMap = map[int]string{
	0: "Info",
	1: "Warning",
	2: "Low",
	3: "Medium",
	4: "High",
	5: "Critical",
}

type Detector struct {
	Name string
	Info Info
	Rule Operator
}

type Info struct {
	Description string
	Severity    int
	CWE         int
}

func (d Detector) EvaluateRule(yamlFilePath string) ([]int, error) {
	ruleType := strings.Split(reflect.TypeOf(d.Rule).String(), ".")
	if slices.Contains(operators, ruleType[len(ruleType)-1]) {
		if err := d.Rule.Evaluate(yamlFilePath); err != nil {
			return nil, err
		}

		return d.Rule.GetLines(), nil
	}

	var stack []Operator
	var queue []Operator

	queue = append(queue, d.Rule)

	for len(queue) > 0 {
		op := queue[0]
		queue = queue[1:]

		opType := strings.Split(reflect.TypeOf(op).String(), ".")

		if slices.Contains(logicalOperators, opType[len(opType)-1]) {
			children := op.GetChildren()

			stack = append(stack, op)
			queue = append(queue, children[0])
			queue = append(queue, children[1])
		}
	}

	for len(stack) > 0 {
		op := stack[len(stack)-1]

		if err := op.Evaluate(yamlFilePath); err != nil {
			return []int{}, err
		}

		stack = stack[:len(stack)-1]
	}

	return d.Rule.GetLines(), nil
}
