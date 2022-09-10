package worker

import "farmer/internal/pkg/entities"

type IWavetrendWorker interface {
	GetCurrentTci() float64
	GetCurrentDifWavetrend() float64
	GetClosePrice() float64
	GetPastWaveTrendData() (*entities.PastWavetrend, bool)

	Run(done chan<- error)
	Stop()
}
