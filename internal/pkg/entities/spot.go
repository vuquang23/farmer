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
	Qty                 float64
	CummulativeQuoteQty float64
	Ref                 uint64
	IsDone              bool
	UnitBought          uint64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
