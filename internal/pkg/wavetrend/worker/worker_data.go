package worker

import (
	en "farmer/internal/pkg/entities"
)

type pastWavetrendData struct {
	lastOpenTime uint64
	lastD        float64
	lastEsa      float64
	pastTci      []float64
	difWavetrend []float64
}

func newPastWavetrendData() *pastWavetrendData {
	return &pastWavetrendData{}
}

func (w *worker) loadLastOpenTime() uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.pastData.lastOpenTime
}

func (w *worker) storePastWaveTrendData(pastData en.PastWavetrend) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.pastData.lastOpenTime = pastData.LastOpenTime
	w.pastData.lastD = pastData.LastD
	w.pastData.lastEsa = pastData.LastEsa
	w.pastData.pastTci = pastData.PastTci
	w.pastData.difWavetrend = pastData.DifWavetrend
}

func (w *worker) loadPastWaveTrendData() en.PastWavetrend {
	w.mu.Lock()
	defer w.mu.Unlock()

	return en.PastWavetrend{
		LastOpenTime: w.pastData.lastOpenTime,
		LastD:        w.pastData.lastD,
		LastEsa:      w.pastData.lastEsa,
		PastTci:      w.pastData.pastTci,
		DifWavetrend: w.pastData.difWavetrend,
	}
}

type currentWavetrendData struct {
	currentTci          float64
	currentDifWavetrend float64 // tci - average(tci, 4)
	closePrice          float64
}

func newCurrentWavetrendData() *currentWavetrendData {
	return &currentWavetrendData{
		currentTci:          0,
		currentDifWavetrend: 0,
		closePrice:          0,
	}
}

func (w *worker) storeCurrentTci(currentTci float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.currentData.currentTci = currentTci
}

func (w *worker) loadCurrentTci() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.currentData.currentTci
}

func (w *worker) storeCurrentDifWavetrend(value float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.currentData.currentDifWavetrend = value
}

func (w *worker) loadCurrentDifWavetrend() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.currentData.currentDifWavetrend
}

func (w *worker) loadClosePrice() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.currentData.closePrice
}

func (w *worker) storeClosePrice(value float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.currentData.closePrice = value
}
