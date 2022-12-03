package entities

import "time"

type HistorySpotTrade struct {
	ID                  uint64 `gorm:"primaryKey"`
	Symbol              string
	Side                string
	BinanceOrderID      uint64
	Qty                 string
	CummulativeQuoteQty string
	Price               float64
	Ref                 uint64
	UnitBought          uint64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewHistorySpotTrade(spotTrade SpotTrade) HistorySpotTrade {
	return HistorySpotTrade{
		ID:                  spotTrade.ID,
		Symbol:              spotTrade.Symbol,
		Side:                spotTrade.Side,
		BinanceOrderID:      spotTrade.BinanceOrderID,
		Qty:                 spotTrade.Qty,
		CummulativeQuoteQty: spotTrade.CummulativeQuoteQty,
		Price:               spotTrade.Price,
		Ref:                 spotTrade.Ref,
		UnitBought:          spotTrade.UnitBought,
	}
}
