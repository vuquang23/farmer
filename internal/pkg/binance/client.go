package binance

import (
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"

	"farmer/internal/pkg/config"
)

var spot *binance.Client

func InitBinanceSpotClient(isTest bool) {
	if spot == nil {
		if isTest {
			binance.UseTestnet = true
		}
		cfg := config.Instance().Binance
		spot = binance.NewClient(cfg.ApiKey, cfg.SecretKey)
	}
}

func BinanceSpotClientInstance() *binance.Client {
	return spot
}

var future *futures.Client

func InitBinanceFutureClient() {
	if future == nil {
		cfg := config.Instance().Binance
		future = binance.NewFuturesClient(cfg.ApiKey, cfg.SecretKey)
	}
}

func BinanceFutureClientInstance() *futures.Client {
	return future
}
