package detector

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"tool/app/internal/helpers"
)

var SeverityMap = map[int]string{
	0: "Info",
	1: "Warning",
	2: "Low",
	3: "Medium",
	4: "High",
	5: "Critical",
}

var ColorMap = map[int]string{
	0: "\033[34;1m",
	1: "\033[37;1m",
	2: "\033[34;1m",
	3: "\033[33;1m",
	4: "\033[31;1m",
	5: "\033[35;1m",
}

type Detector struct {
	Name string
	Info Info
	Rule Operator
}

type Info struct {
	Description string
	Message     string
	Severity    int
	CWE         int
}

func (d Detector) GetSeverity() string {
	return SeverityMap[d.Info.Severity]
}

func (d Detector) EvaluateRule(yamlContent []byte) ([]int, error) {
	ruleType := strings.Split(reflect.TypeOf(d.Rule).String(), ".")

	if slices.Contains(operators, ruleType[len(ruleType)-1]) {
		if err := d.Rule.Evaluate(yamlContent); err != nil {
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

		if err := op.Evaluate(yamlContent); err != nil {
			return []int{}, err
		}

		stack = stack[:len(stack)-1]
	}

	return d.Rule.GetLines(), nil
}

func (d Detector) PrintResults(yamlContent []byte) {
	results := d.Rule.GetLines()

	if len(results) == 0 {
		return
	}

	fmt.Print(
		"\033[39;1m"+strings.ToTitle(d.Name), ColorMap[d.Info.Severity],
		" ["+strings.ToTitle(SeverityMap[d.Info.Severity])+"]", "\033[0m\n\n",
	)

	for _, result := range results {
		line, err := helpers.ReadLine(strings.NewReader(string(yamlContent)), result)

		if err != nil {
			return
		}

		fmt.Println(result, strings.TrimSpace(line))
		fmt.Print(
			strings.Repeat(" ", len(strconv.Itoa(result))),
			" "+ColorMap[d.Info.Severity]+"^~~~ "+d.Info.Message+"\033[0m\n",
		)
	}

	fmt.Println()
}
