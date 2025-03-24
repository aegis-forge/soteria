package commands

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

		if err = parseAndAnalyze(wd+"/.github/workflows", &stats, flags, lines, detects, config); err != nil {
			return err
		}
	}

	aggregated := statistics.AggStatistics{}
	aggregated.Init(flags.Check.Repo, flags.Check.Output)
	aggregated.Aggregate(stats)

	splitName := strings.Split(aggregated.WorkflowName, "/")

	if flags.Check.String {
		if contents, err := json.Marshal(aggregated); err != nil {
			return err
		} else {
			fmt.Print(string(contents[:]))
		}
	} else {
		if err = aggregated.Detectors.SaveToFile(flags.Check.Output, splitName[len(splitName)-1]); err != nil {
			return err
		}

		statistics.GenerateTableDetectors(stats, flags.Check.MaxRows)
		statistics.GenerateAggregatedTableDetectors(aggregated)
	}

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

	if err != nil && !flags.Check.String {
		return err
	}

	if fileInfo != nil && fileInfo.IsDir() {
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
		var yamlContent []byte

		if flags.Check.String {
			yamlContent = []byte(path)
		} else {
			yamlContent, err = os.ReadFile(path)

			if err != nil {
				return err
			}
		}

		linesWorkflow, err := detects.EvaluateWorkflow(path, yamlContent, flags.Check.Verbose)

		if err != nil {
			return err
		}

		if _, ok := lines[path]; !ok {
			lines[path] = linesWorkflow
		}

		var filename string

		if flags.Check.String {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			filename = "input/" + strconv.Itoa(r.Intn(1000))
		} else {
			filename = path
		}

		stat := statistics.Statistics{WorkflowName: filename}
		stat.Init()

		err = stat.ComputeDetectors(yamlContent, lines[path], filename, detects)

		if err != nil {
			return err
		}

		if !flags.Check.Global {
			noExtName := strings.TrimSuffix(stat.WorkflowName, filepath.Ext(stat.WorkflowName))
			splitName := strings.Split(noExtName, "/")
			err = stat.Detectors.SaveToFile(flags.Check.Output, splitName[len(splitName)-1])

			if err != nil {
				return err
			}
		}

		*stats = append(*stats, stat)
	}

	return nil
}
