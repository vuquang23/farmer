package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/utils/logger"
)

func InitWavetrendCalculatorComponents() {
	logger.InitLogger()

	binance.InitBinanceClient()

	services.InitWaveTrendMomentumService(binance.BinanceClientInstance())
}
