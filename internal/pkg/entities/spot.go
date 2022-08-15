package entities

import "time"

type SpotWorker struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement"`
	Symbol          string
	BuyCountAllowed uint64
	BuyCount        int64
	BuyNotional     float64
}

type SpotTrade struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement"`
	BinanceOrderID uint64
	SpotWorkerID   uint64
	OpenQty        float64
	OpenPrice      float64
	CloseQty       float64
	ClosePrice     float64
	CloseCount     uint64
	IsDone         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
