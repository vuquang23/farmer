package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/utils/logger"
)

func InitSymlistUpdaterComponents() error {
	logger.InitLogger()

	if err := binance.InitBinanceSpotClient(false); err != nil {
		return err
	}

	return nil
}
