package commands

import (
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"tool/app/internal/models"
	"tool/app/internal/statistics"
)

func Stats(ctx *cli.Context, flags models.Flags) error {
	var stats []statistics.Statistics

	if ctx.Args().Present() {
		for ind := range ctx.Args().Len() {
			if err := parseAndCompute(ctx.Args().Get(ind), &stats, flags); err != nil {
				return err
			}
		}
	} else {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		if err = parseAndCompute(wd+".github/workflows", &stats, flags); err != nil {
			return err
		}
	}

	aggregated := statistics.AggStatistics{}
	aggregated.Init()
	aggregated.Aggregate(stats)

	err := aggregated.Structure.SaveToFile(flags.Stats.Output, aggregated.WorkflowName)

	if err != nil {
		return err
	}

	statistics.GenerateTableStructure(stats, flags.Stats.MaxRows)
	statistics.GenerateAggregatedTableStructure(aggregated)

	return nil
}

func parseAndCompute(path string, stats *[]statistics.Statistics, flags models.Flags) error {
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
				if err = parseAndCompute(path, stats, flags); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	} else {
		stat := statistics.Statistics{WorkflowName: path}
		stat.Init()
		yamlContent, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		err = stat.ComputeStructure(yamlContent)

		if err != nil {
			return err
		}

		err = stat.Structure.SaveToFile(flags.Stats.Output, stat.WorkflowName)

		*stats = append(*stats, stat)
	}

	return nil
}
