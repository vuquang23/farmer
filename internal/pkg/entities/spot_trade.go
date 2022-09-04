package entities

import "time"

type (
	SpotWorker struct {
		ID             uint64 `gorm:"primaryKey;autoIncrement"`
		Symbol         string // pair
		UnitBuyAllowed uint64
		UnitNotional   float64
	}

	SpotWorkerStatus struct {
		SpotWorker
		TotalUnitBought uint64
	}
)

type SpotTrade struct {
	ID                  uint64 `gorm:"primaryKey;autoIncrement"`
	Side                string
	BinanceOrderID      uint64
	SpotWorkerID        uint64
	Qty                 string
	CummulativeQuoteQty float64
	Price               float64
	Ref                 uint64
	IsDone              bool
	UnitBought          uint64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// Entities for serving telebot
type (
	TradingPairInfo struct {
		Symbol                 string
		UsdBenefit             float64
		BaseAmount             float64
		QuoteAmount            float64
		CurrentUsdValue        float64
		CurrentUsdValueChanged float64
		UnitBuyAllowed         uint64
		UnitNotional           float64
		TotalUnitBought        uint64
	}
)
