package entities

// Entities for serving telebot
type (
	SpotTradingPairInfo struct {
		Symbol          string
		Capital         float64
		CurrentUSDValue float64
		BenefitUSD      float64
		BaseAmount      float64 // base asset current balance.
		QuoteAmount     float64 // quote asset current balance.
		UnitBuyAllowed  uint64
		UnitNotional    float64 // notional ($) of each unit
		TotalUnitBought uint64
	}
)
