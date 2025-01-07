package statistics

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"path/filepath"
	"strings"
	"tool/app/internal/detectors"
)

// ====================
// ==== STATISTICS ====
// ====================

type Statistics struct {
	WorkflowName string    `json:"workflow"`
	Structure    Structure `json:"structure"`
	Detectors    Detectors `json:"detectors"`
}

func (s *Statistics) Init() {
	s.Structure.Workflow = map[string]Group{}
	s.Structure.Jobs = map[string]Group{}
	s.Structure.Steps = map[string]Group{}
	s.Structure.Containers = map[string]Group{}

	s.Detectors.Frequencies = map[string]Group{}
	s.Detectors.Severities = map[string]Group{}

	path := strings.TrimSuffix(s.WorkflowName, filepath.Ext(s.WorkflowName))
	filenameArr := strings.Split(path, "/")
	s.WorkflowName = filenameArr[len(filenameArr)-1]
}

func (s *Statistics) Compute(yamlContent []byte, lines map[string][]int, detects detectors.Detectors) error {
	if err := s.Structure.Compute(yamlContent, s.WorkflowName); err != nil {
		return err
	}

	if err := s.Detectors.Compute(yamlContent, lines, detects); err != nil {
		return err
	}

	return nil
}

func (s *Statistics) SaveToFile() error {
	contents, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		return err
	}

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	fullPath := filepath.Join(wd + "/out/" + s.WorkflowName + ".json")

	if err = os.WriteFile(fullPath, contents, 0644); err != nil {
		return err
	}

	return nil
}

func GenerateTables(statistics []Statistics, maxRows int) {
	fmt.Println()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Statistics per Workflow – Structure")
	t.AppendHeader(table.Row{"NAME", "JOBS", "STEPS", "CONTAINERS"})
	t.AppendRows(createRows(statistics, maxRows))
	t.SetIndexColumn(1)
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Render()

	if len(statistics) > maxRows {
		fmt.Println("...only showing first ", maxRows, " rows...")
	}

	fmt.Println()

	td := table.NewWriter()
	td.SetOutputMirror(os.Stdout)
	td.SetTitle("Statistics per Workflow – Detectors")
	td.AppendHeader(table.Row{"", "INFO", "WARN", "LOW", "MED", "HIGH", "CRIT"})
	td.AppendRows(createRowsDetectors(statistics, maxRows))
	td.SetIndexColumn(1)
	td.SetStyle(table.StyleColoredRedWhiteOnBlack)
	td.Render()

	if len(statistics) > maxRows {
		fmt.Println("...only showing first ", maxRows, " rows...")
	}

	fmt.Println()
}

func createRows(statistics []Statistics, maxRows int) []table.Row {
	var rows []table.Row

	for ind, stat := range statistics {
		if ind >= maxRows {
			break
		}

		rows = append(rows, table.Row{
			stat.WorkflowName,
			stat.Structure.Workflow["jobs"].Frequencies,
			stat.Structure.Jobs["steps"].Frequencies,
			stat.Structure.Jobs["containers"].Frequencies,
		})
	}

	return rows
}

func createRowsDetectors(stats []Statistics, maxRows int) []table.Row {
	var rows []table.Row

	for ind, stat := range stats {
		if ind >= maxRows {
			break
		}

		row := table.Row{}

		row = append(row, stat.WorkflowName)

		for _, severity := range SeveritiesNames {
			if el, ok := stat.Detectors.Severities[severity]; ok {
				row = append(row, el.Frequencies)
			} else {
				row = append(row, 0)
			}
		}

		rows = append(rows, row)
	}

	return rows
}

// ===================
// ==== STRUCTURE ====
// ===================

type Structure struct {
	Workflow   map[string]Group `json:"workflows"`
	Jobs       map[string]Group `json:"jobs"`
	Steps      map[string]Group `json:"steps"`
	Containers map[string]Group `json:"containers"`
}

func (s *Structure) Compute(yamlContent []byte, workflowName string) error {
	if workflow, err := computeWorkflows(yamlContent, workflowName); err == nil {
		s.Workflow = workflow
	} else {
		return err
	}

	if jobs, err := computeJobs(yamlContent); err == nil {
		s.Jobs = jobs
	} else {
		return err
	}

	return nil
}

// ===================
// ==== DETECTORS ====
// ===================

type Detectors struct {
	Severities  map[string]Group `json:"severities"`
	Frequencies map[string]Group `json:"frequencies"`
}

func (d *Detectors) Compute(yamlContent []byte, lines map[string][]int, detects detectors.Detectors) error {
	if severities, frequencies, err := computeDetectors(yamlContent, lines, detects); err == nil {
		d.Severities = severities
		d.Frequencies = frequencies
	} else {
		return err
	}

	return nil
}
