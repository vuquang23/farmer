package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/config"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
	"farmer/internal/pkg/utils/logger"
)

func InitSpotFarmerComponents(isTest bool) {
	logger.InitLogger()

	telebot.InitTeleBot(config.Instance().Telebot.Token, int64(config.Instance().Telebot.GroupID))

	binance.InitBinanceSpotClient(isTest)

	spotmanager.InitSpotManager(binance.BinanceSpotClientInstance())
}
