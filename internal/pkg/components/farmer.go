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

func InitSpotFarmerComponents(isTest bool) {
	// logger
	logger.InitLogger()

	// binance client
	binance.InitBinanceSpotClient(isTest)

	// db
	if err := db.InitDB(); err != nil {
		panic(err)
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

	// telebot
	telebot.InitTeleBot(services.SpotTradeServiceInstance())

	// wavetrend provider
	wtp.InitWavetrendProvider()

	// trading system
	spotmanager.InitSpotManager(
		binance.BinanceSpotClientInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)
}
