package statistics

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"path/filepath"
)

var SeveritiesNames = []string{"info", "warning", "low", "medium", "high", "critical"}

// =======================
// ==== AggSTATISTICS ====
// =======================

type AggStatistics struct {
	WorkflowName string       `json:"workflow"`
	Structure    AggStructure `json:"structure"`
	Detectors    AggDetectors `json:"detectors"`
}

func (a *AggStatistics) Init(repository string) {
	a.WorkflowName = "global"

	if repository != "" {
		a.WorkflowName = repository
	}

	a.Structure.Workflow = map[string]AggGroup{}
	a.Structure.Jobs = map[string]AggGroup{}
	a.Structure.Steps = map[string]AggGroup{}
	a.Structure.Containers = map[string]AggGroup{}

	a.Detectors.Frequencies = map[string]AggGroup{}
	a.Detectors.Severities = map[string]AggGroup{}
}

func (a *AggStatistics) Aggregate(stats []Statistics) {
	aggregated := Statistics{WorkflowName: a.WorkflowName}
	aggregated.Init()

	for _, stat := range stats {
		aggregateStats(stat.Structure.Workflow, a.Structure.Workflow)
		aggregateStats(stat.Structure.Jobs, a.Structure.Jobs)
		aggregateStats(stat.Detectors.Frequencies, a.Detectors.Frequencies)
		aggregateStats(stat.Detectors.Severities, a.Detectors.Severities)
	}
}

func aggregateStats(toAggregate map[string]Group, aggregated map[string]AggGroup) {
	for stat, group := range toAggregate {
		var newOccurrences []string

		if group.Workflow != "" {
			newOccurrences = append(newOccurrences, group.Workflow)
		}
		
		newOccurrences = append(newOccurrences, group.Occurrences...)

		if el, ok := aggregated[stat]; !ok {
			aggregated[stat] = AggGroup{
				Occurrences: [][]string{newOccurrences},
				Frequencies: []int{group.Frequencies},
			}
		} else {
			el.Append(newOccurrences, group.Frequencies)
			aggregated[stat] = el
		}
	}
}

func GenerateAggregatedTableStructure(aggregated AggStatistics) {
	count := aggregated.Structure.Workflow["count"]
	jobs := aggregated.Structure.Workflow["jobs"]
	steps := aggregated.Structure.Jobs["steps"]
	cont := aggregated.Structure.Jobs["containers"]

	ts := table.NewWriter()
	ts.SetOutputMirror(os.Stdout)
	ts.SetTitle("Aggregated Statistics – Structure")
	ts.AppendHeader(table.Row{"", "COUNT", "MIN", "MAX", "MEAN", "MEDIAN", "STD"})
	ts.AppendRow(table.Row{"WORKFLOWS", count.Count()})
	ts.AppendRows([]table.Row{
		{"JOBS", jobs.Count(), jobs.Min(), jobs.Max(), jobs.Mean(), jobs.Median(), jobs.StdDev()},
		{"STEPS", steps.Count(), steps.Min(), steps.Max(), steps.Mean(), steps.Median(), steps.StdDev()},
		{"CONTAINERS", cont.Count(), cont.Min(), cont.Max(), cont.Mean(), cont.Median(), cont.StdDev()},
	})
	ts.SetIndexColumn(1)
	ts.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	ts.Render()

	fmt.Println()
}

func GenerateAggregatedTableDetectors(aggregated AggStatistics) {
	td := table.NewWriter()
	td.SetOutputMirror(os.Stdout)
	td.SetTitle("Aggregated Statistics – Detectors")
	td.AppendHeader(table.Row{"", "COUNT", "MIN", "MAX", "MEAN", "MEDIAN", "STD"})
	td.AppendRows(createAggRows(aggregated.Detectors.Frequencies, false))
	td.AppendSeparator()
	td.AppendRows(createAggRows(aggregated.Detectors.Severities, true))
	td.SetIndexColumn(1)
	td.SetStyle(table.StyleColoredRedWhiteOnBlack)
	td.Render()
}

func createAggRows(stats map[string]AggGroup, severities bool) []table.Row {
	var row []table.Row

	if severities {
		for _, key := range SeveritiesNames {
			stat := stats[key]

			if stat.Count() > 0 {
				row = append(row, table.Row{
					key, stat.Count(), stat.Min(), stat.Max(), stat.Mean(), stat.Median(), stat.StdDev(),
				})
			} else {
				row = append(row, table.Row{key, 0, 0, 0, 0, 0, 0})
			}
		}
	} else {
		for key, stat := range stats {
			row = append(row, table.Row{
				key, stat.Count(), stat.Min(), stat.Max(), stat.Mean(), stat.Median(), stat.StdDev(),
			})
		}
	}

	return row
}

// ======================
// ==== AggSTRUCTURE ====
// ======================

type AggStructure struct {
	Workflow   map[string]AggGroup `json:"workflows"`
	Jobs       map[string]AggGroup `json:"jobs"`
	Steps      map[string]AggGroup `json:"steps"`
	Containers map[string]AggGroup `json:"containers"`
}

func (a *AggStructure) SaveToFile(path string, filename string) error {
	contents, err := json.Marshal(a)

	if err != nil {
		return err
	}

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	fullPath := filepath.Join(wd + "/out/stats/" + filename + ".json")

	if path != "" {
		fullPath = path + "/" + filename + ".json"
	} else {
		err = os.MkdirAll(wd+"/out/stats/", 0755)

		if err != nil {
			return err
		}
	}

	if err = os.WriteFile(fullPath, contents, 0644); err != nil {
		return err
	}

	return nil
}

// ======================
// ==== AggDETECTORS ====
// ======================

type AggDetectors struct {
	Severities  map[string]AggGroup `json:"severities"`
	Frequencies map[string]AggGroup `json:"frequencies"`
}

func (a *AggDetectors) SaveToFile(path string, filename string) error {
	contents, err := json.Marshal(a)

	if err != nil {
		return err
	}

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	fullPath := filepath.Join(wd + "/out/results/" + filename + ".json")

	if path != "" {
		fullPath = path + "/" + filename + ".json"
	} else {
		err = os.MkdirAll(wd+"/out/results/", 0755)

		if err != nil {
			return err
		}
	}

	if err = os.WriteFile(fullPath, contents, 0644); err != nil {
		return err
	}

	return nil
}
