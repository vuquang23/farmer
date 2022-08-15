package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/config"
	"farmer/internal/pkg/db"
	"farmer/internal/pkg/repositories"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
	"farmer/internal/pkg/utils/logger"
)

func InitSpotFarmerComponents(isTest bool) {
	logger.InitLogger()

	telebot.InitTeleBot(config.Instance().Telebot.Token, int64(config.Instance().Telebot.GroupID))

	binance.InitBinanceSpotClient(isTest)

	db.InitDB()

	// repo
	repositories.InitSpotWorkerRepository(db.Instance())

	spotmanager.InitSpotManager(
		binance.BinanceSpotClientInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)
}
