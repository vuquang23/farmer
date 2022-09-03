package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/db"
	"farmer/internal/pkg/repositories"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
	"farmer/internal/pkg/utils/logger"
	wtp "farmer/internal/pkg/wavetrend"
)

func InitSpotFarmerComponents(isTest bool) {
	logger.InitLogger()

	telebot.InitTeleBot()

	binance.InitBinanceSpotClient(isTest)

	if err := db.InitDB(); err != nil {
		panic(err)
	}

	// repo
	repositories.InitSpotWorkerRepository(db.Instance())
	repositories.InitSpotTradeRepository(db.Instance())

	wtp.InitWavetrendProvider()

	spotmanager.InitSpotManager(
		binance.BinanceSpotClientInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)
}
