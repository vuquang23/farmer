package entities

import "github.com/adshao/go-binance/v2"

type (
	SpotBuyOrder struct {
		UnitBought int64
		Price      float64
	}

	SpotBuySignal struct {
		ShouldBuy bool
		Order     SpotBuyOrder
	}
)

type (
	SpotSellOrder struct {
		Qty        string
		UnitBought uint64
		Ref        uint64
		Price      float64
	}

	SpotSellSignal struct {
		ShouldSell bool
		Orders     []*SpotSellOrder
	}

	CreateSpotSellOrderResponse struct {
		BinanceResponse *binance.CreateOrderResponse
		Order           *SpotSellOrder
	}
)
