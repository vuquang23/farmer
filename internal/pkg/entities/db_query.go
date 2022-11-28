package entities

type AggregatedBuyOrders struct {
	TotalBaseAmount          float64
	TotalCummulativeQuoteQty float64
	TotalUnitBought          uint64
}
