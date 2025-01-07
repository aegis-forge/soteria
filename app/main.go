package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"tool/app/internal/commands"
	"tool/app/internal/detectors"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

func main() {
	flags := models.Flags{}

	detects := detectors.Detectors{}
	detects.Init()

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print the version",
	}

	app := &cli.App{
		Name:                   helpers.Constants.Name,
		Usage:                  "Statically analyze GitHub workflow files with custom detectors",
		Version:                helpers.Constants.Version,
		ArgsUsage:              "[file(s) | directory]",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Authors: []*cli.Author{
			{
				Name:  "Edoardo Riggio",
				Email: "edoardo.riggio@usi.ch",
			},
		},
		Action: func(ctx *cli.Context) error {
			return cli.ShowAppHelp(ctx)
		},
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "Perform a static analysis check on the given file(s) or directory",
				Action: func(ctx *cli.Context) error {
					return commands.Check(ctx, flags, detects)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "stats",
						Aliases:     []string{"s"},
						Usage:       "Compute and display the statistics for the passed workflow(s)",
						Destination: &flags.Check.Stats,
					},
					&cli.IntFlag{
						Name:        "max-rows",
						Aliases:     []string{"r"},
						Usage:       "Maximum number of rows to print for the statistics table",
						Value:       10,
						Destination: &flags.Check.MaxRows,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"o"},
						Usage:       "Output directory for the workflows' statistics (one JSON file per workflow will be generated, plus a global one)",
						Destination: &flags.Check.Output,
					},
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Usage:       "Verbose mode",
						Destination: &flags.Check.Verbose,
					},
				},
				UseShortOptionHandling: true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
