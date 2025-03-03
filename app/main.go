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
		Usage:                  "Statically analyze YAML files with custom detectors",
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
					&cli.StringFlag{
						Name:        "repo",
						Aliases:     []string{"r"},
						Usage:       "Name of the repository (will be used to name the global results file)",
						Destination: &flags.Check.Repo,
					},
					&cli.BoolFlag{
						Name:        "string",
						Aliases:     []string{"s"},
						Usage:       "YAML being passed as string",
						Destination: &flags.Check.String,
					},
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Usage:       "Verbose mode",
						Destination: &flags.Check.Verbose,
					},
					&cli.IntFlag{
						Name:        "max-rows",
						Aliases:     []string{"m"},
						Usage:       "Maximum number of rows to print for the statistics table",
						Value:       10,
						Destination: &flags.Check.MaxRows,
					},
					&cli.BoolFlag{
						Name:        "global",
						Aliases:     []string{"g"},
						Usage:       "Output only one global JSON file per repository with the aggregated results",
						Destination: &flags.Check.Global,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"o"},
						Usage:       "Output directory for the workflows' statistics (one JSON file per workflow will be generated, plus a global one)",
						Destination: &flags.Check.Output,
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
					&cli.StringFlag{
						Name:        "repo",
						Aliases:     []string{"r"},
						Usage:       "Name of the repository (will be used to name the global results file)",
						Destination: &flags.Stats.Repo,
					},
					&cli.IntFlag{
						Name:        "max-rows",
						Aliases:     []string{"m"},
						Usage:       "Maximum number of rows to print for the statistics table",
						Value:       10,
						Destination: &flags.Stats.MaxRows,
					},
					&cli.BoolFlag{
						Name:        "global",
						Aliases:     []string{"g"},
						Usage:       "Output only one global JSON file per repository with the aggregated statistics",
						Destination: &flags.Stats.Global,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"o"},
						Usage:       "Output directory for the workflows' statistics (one JSON file per workflow will be generated, plus a global one)",
						Destination: &flags.Stats.Output,
					},
				},
			},
			{
				Name:  "detectors",
				Usage: "Get the list of all the available detectors. If the 'config' flag is set, it returns the detectors enabled by that configuration file",
				Action: func(ctx *cli.Context) error {
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "config",
						Aliases:     []string{"c"},
						Usage:       "Path to the configuration file",
						Destination: &flags.Check.Config,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
