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
					detects := detectors.Detectors{}
					config := models.Config{}

					if err := config.Read(flags.Check.Config); err != nil {
						return err
					}

					detects.Init(config)

					return commands.Check(ctx, flags, detects)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "config",
						Aliases:     []string{"c"},
						Usage:       "Path to the configuration file",
						Destination: &flags.Check.Config,
					},
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Usage:       "Verbose mode",
						Destination: &flags.Check.Verbose,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"o"},
						Usage:       "Output directory for the workflows' statistics (one JSON file per workflow will be generated, plus a global one)",
						Destination: &flags.Check.Output,
					},
					&cli.IntFlag{
						Name:        "max-rows",
						Aliases:     []string{"r"},
						Usage:       "Maximum number of rows to print for the statistics table",
						Value:       10,
						Destination: &flags.Check.MaxRows,
					},
				},
				UseShortOptionHandling: true,
			},
			{
				Name:  "stats",
				Usage: "Compute file structure statistics on the given file(s) or directory (only for GitHub workflow files)",
				Action: func(ctx *cli.Context) error {
					return commands.Stats(ctx, flags)
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "max-rows",
						Aliases:     []string{"r"},
						Usage:       "Maximum number of rows to print for the statistics table",
						Value:       10,
						Destination: &flags.Stats.MaxRows,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"o"},
						Usage:       "Output directory for the workflows' statistics (one JSON file per workflow will be generated, plus a global one)",
						Destination: &flags.Stats.Output,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
