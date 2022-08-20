package spotworker

import (
	"context"
	"fmt"
	"sync"
	"time"

	c "farmer/internal/pkg/constants"
	e "farmer/internal/pkg/entities"
	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/indicators"
	"farmer/internal/pkg/utils/logger"
)

type waveTrendData struct {
	mu *sync.Mutex

	lastOpenTime uint64
	lastD        float64
	lastEsa      float64
	pastTci      []float64 // len 30
	difWavetrend []float64 // len 6

	currentTci          float64
	currentDifWavetrend float64
}

func newWaveTrendData() *waveTrendData {
	return &waveTrendData{
		mu: &sync.Mutex{},
	}
}

func (wt *waveTrendData) loadLastOpenTime() uint64 {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	return wt.lastOpenTime
}

func (wt *waveTrendData) storePastWaveTrendData(pastData e.PastWavetrend) {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	wt.lastOpenTime = pastData.LastOpenTime
	wt.lastD = pastData.LastD
	wt.lastEsa = pastData.LastEsa
	wt.pastTci = pastData.PastTci
	wt.difWavetrend = pastData.DifWavetrend
}

func (wt *waveTrendData) loadPastWaveTrendData() e.PastWavetrend {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	return e.PastWavetrend{
		LastOpenTime: wt.lastOpenTime,
		LastD:        wt.lastD,
		LastEsa:      wt.lastEsa,
		PastTci:      wt.pastTci,
		DifWavetrend: wt.difWavetrend,
	}
}

func (wt *waveTrendData) storeCurrentTci(currentTci float64) {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	wt.currentTci = currentTci
}

func (wt *waveTrendData) loadCurrentTci() float64 {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	return wt.currentTci
}

func (wt *waveTrendData) storeCurrentDifWavetrend(value float64) {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	wt.currentDifWavetrend = value
}

func (wt *waveTrendData) loadCurrentDifWavetrend() float64 {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	return wt.currentDifWavetrend
}

func (w *spotWorker) updateWaveTrendPeriodically(doneC chan<- error) {
	ID := w.setting.loadSymbol()
	log := logger.WithDescription(fmt.Sprintf("%s - update wave trend periodically", ID))

	if err := w.initWaveTrendPastData(); err != nil {
		doneC <- errors.NewDomainErrorInitWavetrendData(err, ID)
		return
	}

	once := &sync.Once{}
	ticker := time.NewTicker(c.ProcessingFrequencyTime)
	oneMinute := uint64(60000)
	for ; !w.getStopSignal(); <-ticker.C {
		// check whether now is new interval
		lastOpenTime := w.waveTrendDat.loadLastOpenTime()
		now := uint64(time.Now().UnixMilli())
		if now-lastOpenTime > oneMinute*2 {
			err := w.updateWaveTrendForNextInterval(
				lastOpenTime+oneMinute,
				(now-lastOpenTime)/oneMinute-1,
			)
			if err != nil {
				log.Sugar().Error(err)
				continue
			}
		}

		// get candle of current interval
		candle, err := w.bclient.NewKlinesService().
			Symbol(ID).
			Interval("1m").
			Limit(1).
			Do(context.Background())
		if err != nil {
			log.Sugar().Error(err)
			continue
		}

		pastWavetrend := w.waveTrendDat.loadPastWaveTrendData()
		currentTci, currentDifWavetrend := indicators.CalculateCurrentTciAndDifWavetrendFromPastWavetrendDatAndCurrentCandle(
			&pastWavetrend, indicators.SpotKlineToMinimalKline(candle),
			c.EmaLenN1, c.EmaLenN2,
		)
		w.waveTrendDat.storeCurrentTci(currentTci)
		w.waveTrendDat.storeCurrentDifWavetrend(currentDifWavetrend)

		once.Do(func() {
			doneC <- nil
		})
	}
}

func (w *spotWorker) initWaveTrendPastData() error {
	symbol := w.setting.loadSymbol()
	interval := "1m"

	candles, err := w.bclient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(int(c.KlineHistoryLen)).
		Do(context.Background())
	if err != nil {
		return err
	}
	candles = candles[:len(candles)-1] // drop last candle

	pastWavetrend, _ := indicators.CalculatePastWavetrendData(
		indicators.SpotKlineToMinimalKline(candles),
		c.EmaLenN1, c.EmaLenN2,
	)

	w.waveTrendDat.storePastWaveTrendData(*pastWavetrend)
	return nil
}

func (w *spotWorker) updateWaveTrendForNextInterval(fromOpenTime uint64, limit uint64) error {
	symbol := w.setting.loadSymbol()
	interval := "1m"

	candles, err := w.bclient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(int64(fromOpenTime)).
		Limit(int(limit)).
		Do(context.Background())
	if err != nil {
		return err
	}

	pastWaveTrend := w.waveTrendDat.loadPastWaveTrendData()
	res, err := indicators.CalculatePastWavetrendDataWithNewCandles(
		&pastWaveTrend,
		indicators.SpotKlineToMinimalKline(candles),
		c.EmaLenN1, c.EmaLenN2,
	)
	if err != nil {
		return err
	}

	w.waveTrendDat.storePastWaveTrendData(*res)
	return nil
}
