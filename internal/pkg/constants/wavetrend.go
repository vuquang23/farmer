package constants

const (
	AvgPeriodLen = 4
	EmaLenN1     = 10
	EmaLenN2     = 21

	M1 = "1m"
	H1 = "1h"

	M1TciLen        = 30
	H1TciLen        = 10 // not used for now but must >= AvgPeriodLen + DifWavetrendLen
	DifWavetrendLen = 6
	KlineHistoryLen = 600
)
