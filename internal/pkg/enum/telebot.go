package enum

type SpotStopType string

const (
	MarketNow  SpotStopType = "MarketNow"
	WaitTarget SpotStopType = "WaitTarget"
)
