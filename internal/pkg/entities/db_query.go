package entities

type AggregatedBuyOrders struct {
	TotalBaseQty    float64
	TotalQuoteQty   float64 // total quote asset qty used to buy base asset.
	TotalUnitBought uint64
}
