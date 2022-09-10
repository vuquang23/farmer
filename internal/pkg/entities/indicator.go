package entities

type (
	MinimalKline struct {
		High     float64
		Low      float64
		Close    float64
		OpenTime uint64
	}

	WavetrendMomentum struct {
		Symbol string // eg: BTC, ETH, ...
		Value  float64
	}

	PastWavetrend struct {
		LastOpenTime uint64 // unix mili
		LastD        float64
		LastEsa      float64
		PastTci      []float64
		DifWavetrend []float64
	}
)
