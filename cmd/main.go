package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"farmer/internal/pkg/builder"
	"farmer/internal/pkg/components"
	"farmer/internal/pkg/config"
	"farmer/internal/pkg/enum"
	"farmer/internal/pkg/migrations"
	"farmer/internal/pkg/services"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
	_ "farmer/pkg/errors"
)

func main() {
	app := &cli.App{
		Name: "My Farmer",
		Commands: []*cli.Command{
			spotFarmerCommand(),
			migrationCommand(),
			updateSymbolListCommand(),
			calcWavetrendMomentumCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func spotFarmerCommand() *cli.Command {
	testFlag := "test"
	cfgFlag := "config"
	defaultCfgFile := "internal/pkg/config/file/default.yaml"

	return &cli.Command{
		Name:  "sfarmer",
		Usage: "Run spot farmer system",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  testFlag,
				Value: true,
			},
			&cli.StringFlag{
				Name:  cfgFlag,
				Value: defaultCfgFile,
			},
		},
		Action: func(ctx *cli.Context) error {
			isTest := ctx.Bool(testFlag)
			fmt.Printf("Run spot farmer with test mode: %t", isTest)

			if err := config.Load(ctx.String(cfgFlag)); err != nil {
				return errors.New("can not load config")
			}

			if err := components.InitSpotFarmerComponents(isTest); err != nil {
				return err
			}

			farmer, err := builder.NewSpotFarmerSystem(spotmanager.SpotManagerInstance(), telebot.TeleBotInstance())
			if err != nil {
				return err
			}

			return farmer.Run()
		},
	}
}

func migrationCommand() *cli.Command {
	cfgFile := "internal/pkg/config/file/default.yaml"
	defaultMigrationDir := "file://./migrations/mysql"
	flagUp := "up"
	flagDown := "down"

	return &cli.Command{
		Name:    "migration",
		Aliases: []string{},
		Usage:   "Run migration",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  flagUp,
				Value: -1,
			},
			&cli.IntFlag{
				Name:  flagDown,
				Value: -1,
			},
		},
		Action: func(c *cli.Context) (err error) {
			err = config.Load(cfgFile)
			if err != nil {
				return err
			}

			up := c.Int(flagUp)
			down := c.Int(flagDown)

			if up == -1 && down == -1 {
				return errors.New("No up or down migration declared")
			}

			if up != -1 && down != -1 {
				return errors.New("[ERROR] Both up and down migration declared. Stop the migration")
			}

			m, err := migrations.NewMigration(defaultMigrationDir)
			if err != nil {
				return err
			}

			if up != -1 {
				return m.MigrateUp(up)
			} else {
				return m.MigrateDown(down)
			}
		},
	}
}

func updateSymbolListCommand() *cli.Command {
	cfgFile := "internal/pkg/config/file/default.yaml"
	return &cli.Command{
		Name:  "symlist",
		Usage: "Update USDT symbol list on Binance",
		Action: func(ctx *cli.Context) error {
			err := config.Load(cfgFile)
			if err != nil {
				return errors.New("can not read config file")
			}

			components.InitSymlistUpdaterComponents()

			updater := builder.NewSymlistUpdater()

			if err := updater.Run(context.Background(), "files/symbol.txt"); err != nil {
				return err
			}
			return nil
		},
	}
}

func calcWavetrendMomentumCommand() *cli.Command {
	cfgFile := "internal/pkg/config/file/default.yaml"
	symlistFile := "files/symbol.txt"
	resultFile := "files/momentum_out.txt"

	symlistFlag := "symlist"
	intervalFlag := "interval"
	marketFlag := "market"

	return &cli.Command{
		Name:  "wtmomentum",
		Usage: "Calculate real (not test market) wavetrend momentum value from a symbol list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  symlistFlag,
				Value: symlistFile,
			},
			&cli.StringFlag{
				Name:  intervalFlag,
				Value: "4h",
			},
			&cli.StringFlag{
				Name:  marketFlag,
				Value: string(enum.Future),
			},
		},
		Action: func(ctx *cli.Context) error {
			err := config.Load(cfgFile)
			if err != nil {
				return err
			}

			components.InitWavetrendCalculatorComponents()

			calculator := builder.NewWaveTrendCalculator(
				services.WaveTrendMomentumServiceInstance(),
			)

			err = calculator.Run(
				context.Background(), enum.Market(ctx.String(marketFlag)),
				ctx.String(intervalFlag), ctx.String(symlistFlag), resultFile,
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
