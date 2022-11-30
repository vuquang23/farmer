package wavetrendprovider

import (
	"context"

	"farmer/internal/pkg/entities"
	errPkg "farmer/pkg/errors"
)

type IWavetrendProvider interface {
	StartService(ctx context.Context, svcName string) *errPkg.DomainError

	SetStopSignal(ctx context.Context, svcName string)

	GetCurrentTci(ctx context.Context, svcName string) (float64, bool)
	GetCurrentDifWavetrend(ctx context.Context, svcName string) (float64, bool)
	GetClosePrice(ctx context.Context, svcName string) (float64, bool)
	//GetPastWaveTrendData return past wavetrend data and a bool indicates that data is outdated or not
	GetPastWaveTrendData(ctx context.Context, svcName string) (*entities.PastWavetrend, bool)
}
