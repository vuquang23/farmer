package entities

type MinimalKline struct {
	High     float64
	Low      float64
	Close    float64
	OpenTime uint64
}

type WavetrendMomentum struct {
	Symbol string // eg: BTC, ETH, ...
	Value  float64
}

type PastWavetrend struct {
	LastOpenTime uint64
	LastD        float64
	LastEsa      float64
	PastTci      []float64
	DifWavetrend []float64
}

type SecondaryPastWavetrend struct {
	LastOpenTime uint64
	LastD        float64
	LastEsa      float64
	PastTci      []float64
}
