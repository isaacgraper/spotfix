package cmd

import (
	"log"
	"os"

	"github.com/isaacgraper/spotfix.git/internal/bot"
	"github.com/isaacgraper/spotfix.git/internal/common/config"
	"github.com/urfave/cli/v2"
)

func Run() error {
	process := bot.NewProcess()

	app := &cli.App{
		Name:  "Bot",
		Usage: "Automate Bot - RPA",
		Commands: []*cli.Command{
			{
				Name:  "exec",
				Usage: "Execute the bot process",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "max",
						Usage: "Maximum number of results to process",
						Value: 100,
					},
					&cli.StringFlag{
						Name:  "hour",
						Usage: "Set the hour for process",
					},
					&cli.StringFlag{
						Name:  "category",
						Usage: "Set the category for process",
					},
					&cli.BoolFlag{
						Name:  "filter",
						Usage: "Enable filtering before the execution",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "filter-category",
						Usage: "Enable with filter to process results based on category",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "adjustment",
						Usage: "Enable adjustment process",
						Value: false,
					},
					&cli.IntFlag{
						Name:  "batch",
						Usage: "Batch size for process without filter",
						Value: 10,
					},
				},
				Action: func(ctx *cli.Context) error {
					config := config.Set(
						ctx.String("hour"),
						ctx.String("category"),
						ctx.Bool("filter"),
						// implement filter category with filter
						ctx.Int("max"),
						ctx.Bool("adjustment"),
						ctx.Int("batch"),
					)

					if err := process.Execute(config); err != nil {
						log.Fatalf("error while trying to start the bot: %v", err)
					}
					return nil
				},
			},
		},
	}
	return app.Run(os.Args)
}
