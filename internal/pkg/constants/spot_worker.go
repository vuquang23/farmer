package constants

import "time"

const (
	ProcessingFrequencyTime          = time.Second
	SecondaryProcessingFrequencyTime = time.Minute

	WavetrendOverbought = 50
	WavetrendOversold   = -53

	OversoldRequiredTime              = 5 // 5 minutes
	OversoldPositiveDifWtRequiredTime = 2 // 2 minutes
	OversoldNegativeDifWtRequiredTime = 3 // 3 minutes
)
