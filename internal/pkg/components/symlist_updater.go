package components

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/utils/logger"
)

func InitSymlistUpdaterComponents() {
	logger.InitLogger()

	binance.InitBinanceClient()
}
