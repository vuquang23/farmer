package worker

type IWavetrendWorker interface {
	SetStopSignal()

	GetCurrentTci() float64
	GetCurrentDifWavetrend() float64

	Run(done chan<- error)
}
