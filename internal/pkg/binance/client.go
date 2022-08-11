package binance

import (
	"farmer/internal/pkg/config"

	"github.com/adshao/go-binance/v2"
)

var client *binance.Client

func InitBinanceClient() {
	if client == nil {
		cfg := config.Instance().Binance
		client = binance.NewClient(cfg.ApiKey, cfg.SecretKey)
	}
}

func BinanceClientInstance() *binance.Client {
	return client
}
