package entities

// Entities for serving telebot
type (
	SpotTradingPairInfo struct {
		Symbol          string
		Capital         float64
		CurrentUSDValue float64
		BenefitUSD      float64
		BaseQty         float64
		QuoteQty        float64
		UnitBuyAllowed  uint64
		UnitNotional    float64 // notional ($) of each unit
		TotalUnitBought uint64
	}
)
