package wavetrendprovider

import (
	"farmer/internal/pkg/entities"
	errPkg "farmer/pkg/errors"
)

type IWavetrendProvider interface {
	StartService(svcName string) *errPkg.DomainError

	SetStopSignal(svcName string)

	GetCurrentTci(svcName string) float64
	GetCurrentDifWavetrend(svcName string) float64
	GetClosePrice(svcName string) float64
	//GetPastWaveTrendData return past wavetrend data and a bool indicates that data is outdated or not
	GetPastWaveTrendData(svcName string) (*entities.PastWavetrend, bool)
}
