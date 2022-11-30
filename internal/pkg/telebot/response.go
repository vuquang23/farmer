package telebot

import "farmer/internal/pkg/entities"

// Get account Info
type (
	SpotPairInfo struct {
		Symbol          string
		Capital         float64
		CurrentUSDValue float64
		BenefitUSD      float64
		ChangedUSD      float64
		BaseAmount      float64
		QuoteAmount     float64
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
		Data       *entities.PastWavetrend
		IsOutdated bool
	}
)
