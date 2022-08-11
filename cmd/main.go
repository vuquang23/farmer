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
	// defaultMigrationDir := "file://./migrations/mysql"
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

			// m, err := migrations.NewMigration(defaultMigrationDir)
			// if err != nil {
			// 	fmt.Println("Can not create migration " + err.Error())
			// }

			// if up != -1 {
			// 	return m.MigrateUp(up)
			// } else {
			// 	return m.MigrateDown(down)
			// }
			return nil
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
			}
			return nil
		},
	}
}
