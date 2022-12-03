package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/utils/logger"
)

func InitWavetrendCalculatorComponents(isTest bool) {
	logger.InitLogger()

	binance.InitBinanceSpotClient(isTest)
	binance.InitBinanceFutureClient()

	services.InitWaveTrendMomentumService(
		binance.BinanceSpotClientInstance(),
		binance.BinanceFutureClientInstance(),
	)
}
