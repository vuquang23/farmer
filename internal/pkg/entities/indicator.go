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
