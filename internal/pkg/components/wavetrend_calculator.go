package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/utils/logger"
)

func InitWavetrendCalculatorComponents(isTest bool) error {
	logger.InitLogger()

	if err := binance.InitBinanceSpotClient(isTest); err != nil {
		return err
	}
	binance.InitBinanceFutureClient()

	services.InitWaveTrendMomentumService(
		binance.BinanceSpotClientInstance(),
		binance.BinanceFutureClientInstance(),
	)

	return nil
}
