package wavetrendprovider

import errPkg "farmer/pkg/errors"

type IWavetrendProvider interface {
	StartService(svcName string) *errPkg.DomainError

	SetStopSignal(svcName string)

	GetCurrentTci(svcName string) float64
	GetCurrentDifWavetrend(svcName string) float64
}
