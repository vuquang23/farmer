package telebot

// Get account Info
type (
	SpotPairInfo struct {
		Symbol                 string
		UsdBenefit             float64
		BaseAmount             float64
		QuoteAmount            float64
		CurrentUsdValue        float64
		CurrentUsdValueChanged float64
		UnitBuyAllowed         uint64
		UnitNotional           float64
		TotalUnitBought        uint64
	}

	GetSpotAccountInfoResponse struct {
		Pairs           []*SpotPairInfo
		TotalUsdBenefit float64
	}
)
