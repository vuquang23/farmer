package telebot

// Get account Info
type (
	SpotPairInfo struct {
		Symbol          string
		Capital         float64
		CurrentUSDValue float64
		BenefitUSD      float64
		ChangedUSD      float64
		BaseAmount      float64 // base asset current balance.
		QuoteAmount     float64 // quote asset current balance.
		UnitBuyAllowed  uint64
		UnitNotional    float64
		TotalUnitBought uint64
	}

	GetSpotAccountInfoResponse struct {
		Pairs           []*SpotPairInfo
		TotalBenefitUSD float64
		TotalChangedUSD float64
	}

	GetWavetrendDataResponse struct {
		PastTci             []float64
		CurrentTci          float64
		DifWavetrend        []float64
		CurrentDifWavetrend float64
		ClosePrice          float64
		IsOutdated          bool
	}
)
