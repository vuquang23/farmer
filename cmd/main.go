package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"farmer/internal/pkg/builder"
	"farmer/internal/pkg/components"
	"farmer/internal/pkg/config"
	"farmer/internal/pkg/enum"
	"farmer/internal/pkg/migrations"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/telebot"
	"farmer/internal/pkg/utils/logger"
)

func main() {
	app := &cli.App{
		Name:     "My Farmer",
		Commands: []*cli.Command{},
	}
	app.Commands = append(app.Commands, telebotCommand())
	app.Commands = append(app.Commands, migrationCommand())
	app.Commands = append(app.Commands, updateSymbolListCommand())
	app.Commands = append(app.Commands, calcWavetrendMomentumCommand())

	app.Run(os.Args)
}

func telebotCommand() *cli.Command {
	cfgFile := "internal/pkg/config/file/default.yaml"

	return &cli.Command{
		Name:  "telebot",
		Usage: "Run telebot",
		Action: func(ctx *cli.Context) error {
			if err := config.Load(cfgFile); err != nil {
				return errors.New("[telegram bot] can not load config")
			}
			components.InitTeleBotComponents()
			botCtx := new(gin.Context)
			logger.BindLoggerToGinNormCtx(botCtx, "telegram-bot")
			if err := telebot.InitTeleBot(
				botCtx,
				config.Instance().Telebot.Token, int64(config.Instance().Telebot.GroupID),
			); err != nil {
				return errors.New("[telegram bot] can not init bot")
			}
			telebot.TeleBotInstance().Run()
			return nil
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
				fmt.Println("No up or down migration declared")
				return nil
			}

			if up != -1 && down != -1 {
				return errors.New("[ERROR] Both up and down migration declared. Stop the migration")
			}

			m, err := migrations.NewMigration(defaultMigrationDir)
			if err != nil {
				fmt.Println("Can not create migration " + err.Error())
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

			updaterCtx := new(gin.Context)
			logger.BindLoggerToGinNormCtx(updaterCtx, "Symlist updater")

			if err := updater.Run(updaterCtx, "files/symbol.txt"); err != nil {
				logger.FromGinCtx(updaterCtx).Error(err.Error())
				return nil
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
		Usage: "Calculate wavetrend momentum value from a symbol list",
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
			calculatorCtx := new(gin.Context)
			logger.BindLoggerToGinNormCtx(calculatorCtx, "Wavetrend calculator")

			err = calculator.Run(
				calculatorCtx, enum.Market(ctx.String(marketFlag)),
				ctx.String(intervalFlag), ctx.String(symlistFlag), resultFile,
			)
			if err != nil {
				logger.FromGinCtx(calculatorCtx).Error(err.Error())
				return nil
			}
			return nil
		},
	}
}
