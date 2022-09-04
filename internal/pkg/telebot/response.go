package telebot

// Get account Info
type (
	PairInfo struct {
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

	GetAccountInfoResponse struct {
		Pairs                       []*PairInfo
		TotalUsdBenefit             float64
		CurrentTotalUsdValueChanged float64
	}
)
