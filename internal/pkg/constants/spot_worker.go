package constants

import "time"

const (
	SleepAfterProcessing = time.Second

	WavetrendOverbought = 50
	WavetrendOversold   = -53

	OversoldRequiredTime                = 5 // 5 minutes
	OversoldPositiveDifWtRequiredTime   = 2 // 2 minutes
	OversoldNegativeDifWtRequiredTime   = 3 // 3 minutes
	OverboughtRequiredTime              = 5 // 5 minutes
	OverboughtNegativeDifWtRequiredTime = 2 // 2 minutes
	OverboughtPositiveDifWtRequiredTime = 3 // 3 minutes

	StopBuyAfterBuy = 2 * time.Minute
)

const (
	UnitBuyOnDowntrend = 1
	UnitBuyOnUpTrend   = 3
)
