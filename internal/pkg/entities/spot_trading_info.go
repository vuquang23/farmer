package entities

// Entities for serving telebot
type (
	SpotTradingPairInfo struct {
		Symbol                 string
		UsdBenefit             float64
		BaseAmount             float64
		QuoteAmount            float64
		CurrentUsdValue        float64
		CurrentUsdValueChanged float64
		UnitBuyAllowed         uint64
		UnitNotional           float64 // notional ($) of each unit
		TotalUnitBought        uint64
	}
)
