package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
	"tool/app/internal/statistics"
)

func main() {
	var outputDir string
	var computeStatistics bool

	app := &cli.App{
		Name:  "Tool",
		Usage: "Statically analyze GitHub workflow files with custom detectors",
		Action: func(ctx *cli.Context) error {
			var filenames []string
			var workflows []models.Workflow
			var stats []models.Statistics

			for ind := range ctx.Args().Len() {
				fileInfo, err := os.Stat(ctx.Args().Get(ind))
				if err != nil {
					return err
				}

				if fileInfo.IsDir() {
					err := filepath.Walk(ctx.Args().Get(ind), func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						if !info.IsDir() {
							workflow, err := helpers.ReadFile(path)
							if err != nil {
								return err
							}

							workflows = append(workflows, workflow)
							filenames = append(filenames, path)

							if computeStatistics {
								stat := statistics.ComputeStatistics(workflow)
								stats = append(stats, stat)
							}
						}

						return nil
					})

					if err != nil {
						return err
					}
				} else {
					workflow, err := helpers.ReadFile(ctx.Args().Get(ind))
					if err != nil {
						return err
					}

					workflows = append(workflows, workflow)
					filenames = append(filenames, ctx.Args().Get(ind))

					if computeStatistics {
						stat := statistics.ComputeStatistics(workflow)
						stats = append(stats, stat)
					}
				}
			}

			if strings.Compare(outputDir, "") != 0 {
				outPathStats, err := os.Stat(outputDir)
				if err != nil || !outPathStats.IsDir() {
					return err
				}

				for ind, stat := range stats {
					filename := strings.Split(filenames[ind], "/")
					filenameLast := filename[len(filename)-1]
					filenameAlt := filenameLast[:len(filenameLast)-len(filepath.Ext(filenameLast))]

					if err := helpers.WriteJSONToFile(outputDir+"/"+filenameAlt+"-stats.json", stat); err != nil {
						return err
					}
				}
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "statistics",
				Aliases:     []string{"s", "stats"},
				Usage:       "Compute and display the statistics for the passed workflow(s)",
				Destination: &computeStatistics,
			},
			&cli.StringFlag{
				Name:        "out-dir",
				Aliases:     []string{"o", "out"},
				Usage:       "Output directory for the workflow statistics (one JSON file per workflow will be generated)",
				Destination: &outputDir,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
