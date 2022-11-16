package worker

import (
	"context"

	"farmer/internal/pkg/entities"
)

type IWavetrendWorker interface {
	GetCurrentTci() (float64, bool)
	GetCurrentDifWavetrend() (float64, bool)
	GetClosePrice() (float64, bool)
	GetPastWaveTrendData() (*entities.PastWavetrend, bool)

	Run(ctx context.Context, done chan<- error)
	Stop()
}
