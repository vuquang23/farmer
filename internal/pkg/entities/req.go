package entities

type CreateNewSpotWorkerParams struct {
	Symbol         string
	UnitBuyAllowed uint64
	UnitNotional   float64
}

type StopWorkerParams struct {
	Symbol string
}

type AddCapitalParams struct {
	Symbol  string
	Capital float64
}

type ArchiveSpotTradingDataParams struct {
	Symbol string
}
