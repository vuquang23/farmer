package entities

import "github.com/adshao/go-binance/v2"

type (
	BuyOrder struct {
		UnitBought int64
	}

	BuySignal struct {
		ShouldBuy bool
		Order     BuyOrder
	}
)

type (
	SellOrder struct {
		Qty        string
		UnitBought uint64
		Ref        uint64
	}

	SellSignal struct {
		ShouldSell bool
		Orders     []*SellOrder
	}

	CreateSellOrderResponse struct {
		BinanceResponse *binance.CreateOrderResponse
		Order           *SellOrder
	}
)
