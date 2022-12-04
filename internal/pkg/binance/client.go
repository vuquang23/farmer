package binance

import (
	goctx "context"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/utils/context"
	"farmer/internal/pkg/utils/logger"
)

var spot *binance.Client

func InitBinanceSpotClient(isTest bool) {
	if spot == nil {
		if isTest {
			binance.UseTestnet = true
		}
		cfg := config.Instance().Binance
		spot = binance.NewClient(cfg.ApiKey, cfg.SecretKey)

		spotNotifyCurrentUSDTBalance()
	}
}

func spotNotifyCurrentUSDTBalance() {
	ctx := context.Child(goctx.Background(), "[spot-notify-current-usdt-balance]")
	data, err := spot.NewGetAccountService().Do(ctx)
	if err != nil {
		panic(err)
	}
	for _, d := range data.Balances {
		if d.Asset == "USDT" {
			logger.Infof(ctx, "%+v", d)
			break
		}
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
