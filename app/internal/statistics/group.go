package statistics

import (
	"gonum.org/v1/gonum/stat"
	"slices"
	"tool/app/internal/detector"
	"tool/app/internal/helpers"
)

// ===============
// ==== GROUP ====
// ===============

type Group struct {
	Workflow    string   `json:"workflow,omitempty"`
	Occurrences []string `json:"occurrences"`
	Frequencies int      `json:"frequencies"`
}

func (g *Group) GetOccurrences(yamlPath string, yamlContent []byte) ([]string, int, error) {
	var occurrences []string
	var frequencies int

	res, _, err := detector.Resolve(yamlPath, yamlContent)

	if err != nil {
		return nil, 0, err
	}

	for _, el := range res {
		if el != "" {
			frequencies++
			occurrences = append(occurrences, el)
		}
	}

	return occurrences, frequencies, nil
}

func (g *Group) AddOccurrences(yamlPaths []string, yamlContent []byte) error {
	for _, yamlPath := range yamlPaths {
		occurrences, frequencies, err := g.GetOccurrences(yamlPath, yamlContent)

		if err != nil {
			return err
		}

		if frequencies == 0 {
			continue
		}

		g.Occurrences = append(g.Occurrences, occurrences...)
		g.Frequencies += frequencies
	}

	return nil
}

func (g *Group) AddManually(occurrences []string, frequencies int) {
	g.Occurrences = occurrences
	g.Frequencies = frequencies
}

func CountOccurrences(yamlPath string, yamlContent []byte) int {
	res, _, err := detector.Resolve(yamlPath, yamlContent)

	if err != nil {
		return 0
	}

	return len(res)
}

func checkIfExists(path string, element string, yamlContent []byte) (bool, int, error) {
	res, _, err := detector.Resolve(path, yamlContent)

	if err != nil {
		return false, 0, err
	}

	if element == "*" {
		res1, _, _ := detector.Resolve(path+"[*]", yamlContent)

		return len(helpers.RemoveEmptyStrings(res)) > 0 || len(helpers.RemoveEmptyStrings(res1)) > 0, 1, nil
	} else if slices.Contains(res, element) {
		times := 0

		for _, el := range helpers.RemoveEmptyStrings(res) {
			if el == element {
				times++
			}
		}

		return true, times, nil
	}

	return false, 0, nil
}

// ==================
// ==== AggGROUP ====
// ==================

type AggGroup struct {
	Occurrences [][]string `json:"occurrences"`
	Frequencies []int      `json:"frequencies"`
}

func (g *AggGroup) Append(occurrences []string, frequencies int) {
	if g.Occurrences == nil {
		g.Occurrences = make([][]string, 0)
	}

	g.Occurrences = append(g.Occurrences, occurrences)
	g.Frequencies = append(g.Frequencies, frequencies)
}

func (g *AggGroup) Count() int {
	total := 0

	for _, frequencies := range g.Frequencies {
		total += frequencies
	}

	return total
}

func (g *AggGroup) Min() int {
	if len(g.Frequencies) == 0 {
		return 0
	}

	return slices.Min(g.Frequencies)
}

func (g *AggGroup) Max() int {
	if len(g.Frequencies) == 0 {
		return 0
	}

	return slices.Max(g.Frequencies)
}

func (g *AggGroup) Mean() int {
	if len(g.Frequencies) == 0 {
		return 0
	}

	return g.Count() / len(g.Frequencies)
}

func (g *AggGroup) Median() int {
	if len(g.Frequencies) == 0 {
		return 0
	}

	slices.Sort(g.Frequencies)

	median := len(g.Frequencies) / 2

	if len(g.Frequencies)%2 != 0 {
		return g.Frequencies[median]
	}

	return (g.Frequencies[median-1] + g.Frequencies[median]) / 2
}

func (g *AggGroup) StdDev() float64 {
	if len(g.Frequencies) == 0 {
		return 0
	}

	var floats []float64

	for _, i := range g.Frequencies {
		floats = append(floats, float64(i))
	}

	return float64(int(stat.StdDev(floats, nil)*100)) / 100
}
