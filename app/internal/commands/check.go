package commands

import (
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"tool/app/detectors"
	"tool/app/internal/models"
	"tool/app/internal/statistics"
)

func Check(ctx *cli.Context, flags models.Flags, detects detectors.Detectors) error {
	var stats []statistics.Statistics

	lines := map[string]map[string][]int{}

	if !ctx.Args().Present() {
		if err := parseAndAnalyze(".github/workflows", &stats, flags, lines, detects); err != nil {
			return err
		}
	} else {
		for ind := range ctx.Args().Len() {
			if err := parseAndAnalyze(ctx.Args().Get(ind), &stats, flags, lines, detects); err != nil {
				return err
			}
		}
	}

	if flags.Check.Stats {
		aggregated := statistics.AStatistics{}
		aggregated.Init()
		aggregated.Aggregate(stats)

		err := aggregated.SaveToFile()

		if err != nil {
			return err
		}

		if flags.Check.Verbose {
			statistics.GenerateTables(stats, flags.MaxRows)
			statistics.GenerateAggregatedTables(aggregated)
		}
	}

	return nil
}

func parseAndAnalyze(path string, stats *[]statistics.Statistics, flags models.Flags, lines map[string]map[string][]int, detects detectors.Detectors) error {
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
				if err = parseAndAnalyze(path, stats, flags, lines, detects); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	} else {
		yamlContent, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		linesWorkflow, err := detects.EvaluateWorkflow(path, yamlContent)

		if err != nil {
			return err
		}

		if _, ok := lines[path]; !ok {
			lines[path] = linesWorkflow
		}

		if flags.Check.Stats {
			stat := statistics.Statistics{WorkflowName: path}
			stat.Init()

			err = stat.Compute(yamlContent, lines[path], detects)

			if err != nil {
				return err
			}

			err = stat.SaveToFile()

			if err != nil {
				return err
			}

			*stats = append(*stats, stat)
		}
	}

	return nil
}
