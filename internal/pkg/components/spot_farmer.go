package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/db"
	"farmer/internal/pkg/repositories"
	"farmer/internal/pkg/services"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
	"farmer/internal/pkg/utils/logger"
	wtp "farmer/internal/pkg/wavetrend"
)

func InitSpotFarmerComponents(isTest bool) error {
	// logger
	if err := logger.InitLogger(); err != nil {
		return err
	}

	// binance client
	binance.InitBinanceSpotClient(isTest)

	// db
	if err := db.InitDB(); err != nil {
		return err
	}

	// repo
	repositories.InitSpotWorkerRepository(db.Instance())
	repositories.InitSpotTradeRepository(db.Instance())

	// service
	services.InitSpotTradeService(
		binance.BinanceSpotClientInstance(),
		repositories.SpotTradeRepositoryInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)

	// wavetrend provider
	wtp.InitWavetrendProvider()

	// trading system
	spotmanager.InitSpotManager(
		binance.BinanceSpotClientInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)

	// telebot
	if err := telebot.InitTeleBot(
		services.SpotTradeServiceInstance(),
		wtp.WavetrendProviderInstance(),
		spotmanager.SpotManagerInstance(),
	); err != nil {
		return err
	}

	return nil
}
