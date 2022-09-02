package worker

import "farmer/internal/pkg/entities"

type IWavetrendWorker interface {
	SetStopSignal()

	GetCurrentTci() float64
	GetCurrentDifWavetrend() float64
	GetClosePrice() float64
	GetPastWaveTrendData() *entities.PastWavetrend

	Run(done chan<- error)
}
