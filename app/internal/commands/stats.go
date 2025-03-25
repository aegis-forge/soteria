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

		if err = parseAndCompute(wd+"/.github/workflows", &stats, flags); err != nil {
			return err
		}
	}

	aggregated := statistics.AggStatistics{}
	aggregated.Init(flags.Stats.Repo, flags.Stats.Output)
	aggregated.Aggregate(stats)

	splitName := strings.Split(aggregated.WorkflowName, "/")

	if flags.Stats.String {
		if contents, err := json.Marshal(aggregated); err != nil {
			return err
		} else {
			fmt.Print(string(contents[:]))
		}
	} else {
		if err := aggregated.Structure.SaveToFile(flags.Stats.Output, splitName[len(splitName)-1]); err != nil {
			return err
		}

		statistics.GenerateTableStructure(stats, flags.Stats.MaxRows)
		statistics.GenerateAggregatedTableStructure(aggregated)
	}

	return nil
}

func parseAndCompute(path string, stats *[]statistics.Statistics, flags models.Flags) error {
	fileInfo, err := os.Stat(path)

	if err != nil && !flags.Stats.String {
		return err
	}

	if fileInfo != nil && fileInfo.IsDir() {
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
		var filename string

		if flags.Stats.String {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			filename = "input/" + strconv.Itoa(r.Intn(1000))
		} else {
			filename = path
		}

		stat := statistics.Statistics{WorkflowName: filename}
		stat.Init()

		var yamlContent []byte

		if flags.Stats.String {
			yamlContent = []byte(path)
		} else {
			yamlContent, err = os.ReadFile(path)

			if err != nil {
				return err
			}
		}

		err = stat.ComputeStructure(yamlContent)

		if err != nil {
			return err
		}

		if !flags.Stats.Global && !flags.Stats.String {
			noExtName := strings.TrimSuffix(stat.WorkflowName, filepath.Ext(stat.WorkflowName))
			splitName := strings.Split(noExtName, "/")
			err = stat.Structure.SaveToFile(flags.Stats.Output, splitName[len(splitName)-1])

			if err != nil {
				return err
			}
		}

		*stats = append(*stats, stat)
	}

	return nil
}
