package wavetrendprovider

import (
	"context"

	"farmer/internal/pkg/entities"
	errPkg "farmer/pkg/errors"
)

type IWavetrendProvider interface {
	StartService(ctx context.Context, svcName string) *errPkg.DomainError

	SetStopSignal(svcName string)

	GetCurrentTci(svcName string) (float64, bool)
	GetCurrentDifWavetrend(svcName string) (float64, bool)
	GetClosePrice(svcName string) (float64, bool)
	//GetPastWaveTrendData return past wavetrend data and a bool indicates that data is outdated or not
	GetPastWaveTrendData(svcName string) (*entities.PastWavetrend, bool)
}
