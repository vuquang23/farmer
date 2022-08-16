package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/db"
	"farmer/internal/pkg/repositories"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
	"farmer/internal/pkg/utils/logger"
)

func InitSpotFarmerComponents(isTest bool) {
	logger.InitLogger()

	telebot.InitTeleBot()

	binance.InitBinanceSpotClient(isTest)

	db.InitDB()

	// repo
	repositories.InitSpotWorkerRepository(db.Instance())

	spotmanager.InitSpotManager(
		binance.BinanceSpotClientInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)
}
