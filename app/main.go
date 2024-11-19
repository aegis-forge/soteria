package main

import (
	"flag"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"tool/app/internal/helpers"
)

func main() {
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

							log.Print(workflow)
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

					log.Print(workflow)
				}
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "Compute and display statistics for each passed workflow",
				Action: func(ctx *cli.Context, val bool) error {
					if val {

					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func main_() {
	// CLI flags
	stats := flag.Bool("stats", false, "Compute and display statistics for each passed workflow")
	flag.Usage = func() {
		_, err := fmt.Fprintln(os.Stderr, "Flags:")

		if err != nil {
			return
		}

		flag.PrintDefaults()
	}

	// Parse CLI flags
	flag.Parse()

	log.Print(*stats)
}
