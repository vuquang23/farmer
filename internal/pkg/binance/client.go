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

func InitBinanceSpotClient(isTest bool) error {
	if spot == nil {
		if isTest {
			binance.UseTestnet = true
		}
		cfg := config.Instance().Binance
		spot = binance.NewClient(cfg.ApiKey, cfg.SecretKey)

		if err := spotNotifyCurrentUSDTBalance(); err != nil {
			return err
		}
	}
	return nil
}

func spotNotifyCurrentUSDTBalance() error {
	ctx := context.Child(goctx.Background(), "[spot] notify current USDT balance]")

	data, err := spot.NewGetAccountService().Do(ctx)
	if err != nil {
		return err
	}
	for _, b := range data.Balances {
		if b.Asset == "USDT" {
			logger.Infof(ctx, "%+v", b)
			break
		}
	}

	return nil
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
