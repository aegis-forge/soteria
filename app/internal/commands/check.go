package commands

import (
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"tool/app/internal/detectors"
	"tool/app/internal/models"
	"tool/app/internal/statistics"
)

func Check(ctx *cli.Context, flags models.Flags, detects detectors.Detectors) error {
	var stats []statistics.Statistics

	config, err := readAndValidateConfig(flags.Check.Config)

	if err != nil {
		return err
	}

	lines := map[string]map[string][]int{}

	if ctx.Args().Present() {
		for ind := range ctx.Args().Len() {
			if err = parseAndAnalyze(ctx.Args().Get(ind), &stats, flags, lines, detects, config); err != nil {
				return err
			}
		}
	} else {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		if err = parseAndAnalyze(wd+".github/workflows", &stats, flags, lines, detects, config); err != nil {
			return err
		}
	}

	aggregated := statistics.AggStatistics{}
	aggregated.Init()
	aggregated.Aggregate(stats)

	err = aggregated.Detectors.SaveToFile(flags.Stats.Output, aggregated.WorkflowName)

	if err != nil {
		return err
	}

	statistics.GenerateTableDetectors(stats, flags.Check.MaxRows)
	statistics.GenerateAggregatedTableDetectors(aggregated)

	return nil
}

func readAndValidateConfig(path string) (models.Config, error) {
	var config models.Config

	err := config.Read(path)

	if err != nil {
		return config, err
	}

	return config, err
}

func parseAndAnalyze(path string, stats *[]statistics.Statistics, flags models.Flags, lines map[string]map[string][]int, detects detectors.Detectors, config models.Config) error {
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
				if err = parseAndAnalyze(path, stats, flags, lines, detects, config); err != nil {
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

		linesWorkflow, err := detects.EvaluateWorkflow(path, yamlContent, flags.Check.Verbose)

		if err != nil {
			return err
		}

		if _, ok := lines[path]; !ok {
			lines[path] = linesWorkflow
		}

		stat := statistics.Statistics{WorkflowName: path}
		stat.Init()

		err = stat.ComputeDetectors(yamlContent, lines[path], detects)

		if err != nil {
			return err
		}

		err = stat.Detectors.SaveToFile(flags.Check.Output, stat.WorkflowName)

		if err != nil {
			return err
		}

		*stats = append(*stats, stat)
	}

	return nil
}
