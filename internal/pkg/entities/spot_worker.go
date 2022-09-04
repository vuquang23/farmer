package entities

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
