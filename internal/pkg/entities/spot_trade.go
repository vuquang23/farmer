package entities

import "time"

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
