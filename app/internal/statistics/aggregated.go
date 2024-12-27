package statistics

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"path/filepath"
)

var severitiesNames = []string{"info", "warning", "low", "medium", "high", "critical"}

// =====================
// ==== ASTATISTICS ====
// =====================

type AStatistics struct {
	WorkflowName string     `json:"workflow"`
	Structure    AStructure `json:"structure"`
	Detectors    ADetectors `json:"detectors"`
}

func (a *AStatistics) Init() {
	a.WorkflowName = "global"

	a.Structure.Workflow = map[string]AGroup{}
	a.Structure.Jobs = map[string]AGroup{}
	a.Structure.Steps = map[string]AGroup{}
	a.Structure.Containers = map[string]AGroup{}

	a.Detectors.Frequencies = map[string]AGroup{}
	a.Detectors.Severities = map[string]AGroup{}
}

func (a *AStatistics) Aggregate(stats []Statistics) {
	aggregated := Statistics{WorkflowName: "global"}
	aggregated.Init()

	for _, stat := range stats {
		aggregateStats(stat.Structure.Workflow, a.Structure.Workflow)
		aggregateStats(stat.Structure.Jobs, a.Structure.Jobs)
		aggregateStats(stat.Detectors.Frequencies, a.Detectors.Frequencies)
		aggregateStats(stat.Detectors.Severities, a.Detectors.Severities)
	}
}

func aggregateStats(toAggregate map[string]Group, aggregated map[string]AGroup) {
	for stat, group := range toAggregate {
		if el, ok := aggregated[stat]; !ok {
			aggregated[stat] = AGroup{
				Occurrences: [][]string{group.Occurrences},
				Frequencies: []int{group.Frequencies},
			}
		} else {
			el.Append(group.Occurrences, group.Frequencies)
			aggregated[stat] = el
		}
	}
}

func (a *AStatistics) SaveToFile() error {
	contents, err := json.Marshal(a)

	if err != nil {
		return err
	}

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	fullPath := filepath.Join(wd + "/out/" + a.WorkflowName + ".json")

	if err = os.WriteFile(fullPath, contents, 0644); err != nil {
		return err
	}

	return nil
}

func GenerateAggregatedTables(aggregated AStatistics) {
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

	td := table.NewWriter()
	td.SetOutputMirror(os.Stdout)
	td.SetTitle("Aggregated Statistics – Detectors")
	td.AppendHeader(table.Row{"", "COUNT", "MIN", "MAX", "MEAN", "MEDIAN", "STD"})
	td.AppendRows(createARows(aggregated.Detectors.Frequencies, false))
	td.AppendSeparator()
	td.AppendRows(createARows(aggregated.Detectors.Severities, true))
	td.SetIndexColumn(1)
	td.SetStyle(table.StyleColoredRedWhiteOnBlack)
	td.Render()
}

func createARows(stats map[string]AGroup, severities bool) []table.Row {
	var row []table.Row

	if severities {
		for _, key := range severitiesNames {
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

// ===================
// ==== STRUCTURE ====
// ===================

type AStructure struct {
	Workflow   map[string]AGroup `json:"workflows"`
	Jobs       map[string]AGroup `json:"jobs"`
	Steps      map[string]AGroup `json:"steps"`
	Containers map[string]AGroup `json:"containers"`
}

// ===================
// ==== DETECTORS ====
// ===================

type ADetectors struct {
	Severities  map[string]AGroup `json:"severities"`
	Frequencies map[string]AGroup `json:"frequencies"`
}
