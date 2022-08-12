package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/utils/logger"
)

func InitWavetrendCalculatorComponents() {
	logger.InitLogger()

	binance.InitBinanceSpotClient()
	binance.InitBinanceFutureClient()

	services.InitWaveTrendMomentumService(
		binance.BinanceSpotClientInstance(),
		binance.BinanceFutureClientInstance(),
	)
}
