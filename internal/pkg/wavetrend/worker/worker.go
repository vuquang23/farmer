package worker

import (
	"context"

	"farmer/internal/pkg/entities"
)

type IWavetrendWorker interface {
	GetCurrentTci(ctx context.Context) (float64, bool)
	GetCurrentDifWavetrend(ctx context.Context) (float64, bool)
	GetClosePrice(ctx context.Context) (float64, bool)
	GetPastWaveTrendData(ctx context.Context) (*entities.PastWavetrend, bool)

	Run(ctx context.Context, done chan<- error)
	Stop(ctx context.Context)
}
