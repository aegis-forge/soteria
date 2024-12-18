package commands

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tool/app/detectors"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
	"tool/app/internal/statistics"
)

func Check(ctx *cli.Context, flags models.Flags) error {
	var global models.GlobalStatistics
	var filenames []string
	var workflows []models.Workflow
	var stats []models.Statistics

	lines := map[string]map[string][]int{}

	if !ctx.Args().Present() {
		if err := parseAndAnalyze(".github/workflows", &workflows, &filenames, &stats, flags, lines); err != nil {
			return err
		}
	} else {
		for ind := range ctx.Args().Len() {
			if err := parseAndAnalyze(ctx.Args().Get(ind), &workflows, &filenames, &stats, flags, lines); err != nil {
				return err
			}
		}
	}

	global = statistics.AggregateStatistics(stats)

	if strings.Compare(flags.Check.Output, "") != 0 {
		if err := saveStats(flags, global, stats, filenames); err != nil {
			return err
		}
	}

	if flags.Check.Verbose {
		if err := generateStatsTables(flags, filenames, stats, global); err != nil {
			return err
		}
	}

	log.Print(lines)

	return nil
}

func parseAndAnalyze(path string, workflows *[]models.Workflow, filenames *[]string, stats *[]models.Statistics, flags models.Flags, lines map[string]map[string][]int) error {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				if err = parseAndAnalyze(path, workflows, filenames, stats, flags, lines); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	} else {
		workflow, err := helpers.ReadWorkflow(path)
		if err != nil {
			return err
		}

		*workflows = append(*workflows, workflow)
		*filenames = append(*filenames, path)

		if flags.Check.Stats {
			stat := statistics.ComputeStatistics(workflow)
			*stats = append(*stats, stat)
		}

		detects := detectors.Detectors{}
		detects.Init()

		linesWorkflow, err := detects.EvaluateWorkflow(path)
		if err != nil {
			return err
		}

		if _, ok := lines[path]; !ok {
			lines[path] = linesWorkflow
		}
	}

	return nil
}

func saveStats(flags models.Flags, global models.GlobalStatistics, stats []models.Statistics, filenames []string) error {
	outPathStats, err := os.Stat(flags.Check.Output)
	if err != nil || !outPathStats.IsDir() {
		return err
	}

	for ind, stat := range stats {
		filename := strings.Split(filenames[ind], "/")
		filenameLast := filename[len(filename)-1]
		filenameAlt := filenameLast[:len(filenameLast)-len(filepath.Ext(filenameLast))]

		if err := helpers.WriteJSONToFile(flags.Check.Output+"/"+filenameAlt+"-stats.json", stat); err != nil {
			return err
		}
	}

	if err := helpers.WriteJSONToFile(flags.Check.Output+"/global-stats.json", global); err != nil {
		return err
	}

	return nil
}

func generateStatsTables(flags models.Flags, filenames []string, stats []models.Statistics, global models.GlobalStatistics) error {
	var rows []table.Row

	moreThan := false

	for ind, stat := range stats {
		if ind >= flags.MaxRows {
			moreThan = true
			break
		}

		filename := strings.Split(filenames[ind], "/")

		rows = append(rows, table.Row{
			filename[len(filename)-1],
			stat.Workflow.Jobs.Total,
			stat.Jobs.Steps.Total,
			stat.Jobs.CustomContainers.Total,
		})
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Statistics per Workflow")
	t.AppendHeader(table.Row{"NAME", "JOBS", "STEPS", "CONTAINERS"})
	t.AppendRows(rows)
	t.SetIndexColumn(1)
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Render()

	if moreThan {
		fmt.Println("...only showing first ", flags.MaxRows, " rows...")
	}

	fmt.Println()

	t2 := table.NewWriter()
	t2.SetOutputMirror(os.Stdout)
	t2.SetTitle("Global Statistics")
	t2.AppendHeader(table.Row{"", "COUNT", "MIN", "MAX", "MEAN", "MEDIAN", "STD"})
	t2.AppendRows([]table.Row{
		{"WORKFLOWS", len(filenames)},
		{
			"JOBS", global.Jobs.Total, global.Jobs.Min, global.Jobs.Max, global.Jobs.Mean, global.Jobs.Median,
			float64(int(global.Jobs.Std*10)) / 10,
		},
		{
			"STEPS", global.Steps.Total, global.Steps.Min, global.Steps.Max, global.Steps.Mean, global.Steps.Median,
			float64(int(global.Steps.Std*10)) / 10,
		},
		{
			"CONTAINERS", global.Containers.Total, global.Containers.Min, global.Containers.Max, global.Containers.Mean,
			global.Containers.Median, float64(int(global.Containers.Std*10)) / 10,
		},
	})
	t2.SetIndexColumn(1)
	t2.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t2.Render()

	return nil
}
