package entities

type CreateNewSpotWorkerParams struct {
	Symbol         string
	UnitBuyAllowed uint64
	UnitNotional   float64
}

type StopBotParams struct {
	Symbol string
}