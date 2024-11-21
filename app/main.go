package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"tool/app/internal/helpers"
	"tool/app/internal/statistics"
)

func main() {
	var computeStatistics bool

	app := &cli.App{
		Name:  "Tool",
		Usage: "Statically analyze GitHub workflow files with custom detectors",
		Action: func(ctx *cli.Context) error {
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

							if computeStatistics {
								stats := statistics.ComputeStatistics(workflow)
								log.Print(stats)
							}

							//log.Print(workflow)
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

					if computeStatistics {
						stats := statistics.ComputeStatistics(workflow)
						log.Print(stats.Jobs)
					}

					//log.Print(workflow)
				}
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "stats",
				Aliases:     []string{"s"},
				Usage:       "Compute and display statistics for each passed workflow",
				Destination: &computeStatistics,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
