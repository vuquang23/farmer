package spotworker

import (
	"context"
	e "farmer/internal/pkg/entities"
	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/indicators"
	"farmer/internal/pkg/utils/logger"
	"fmt"
	"sync"
	"time"
)

type waveTrendData struct {
	mu *sync.Mutex

	lastOpenTime uint64
	lastD        float64
	lastEsa      float64
	pastTci      []float64 // len 30

	currentTci float64
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
}

func (wt *waveTrendData) loadPastWaveTrendData() e.PastWavetrend {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	return e.PastWavetrend{
		LastOpenTime: wt.lastOpenTime,
		LastD:        wt.lastD,
		LastEsa:      wt.lastEsa,
		PastTci:      wt.pastTci,
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

func (w *spotWorker) updateWaveTrendPeriodically(doneC chan<- error) {
	ID := w.setting.loadSymbol()
	log := logger.WithDescription(fmt.Sprintf("%s - update wave trend periodically", ID))

	if err := w.initWaveTrendPastData(); err != nil {
		doneC <- errors.NewDomainErrorInitWavetrendData(err, ID)
		return
	}

	once := &sync.Once{}
	ticker := time.NewTicker(time.Second)
	oneMinute := uint64(60000)
	for ; true; <-ticker.C {
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
		currentTci := indicators.CalculateCurrentTciFromPastWavetrendDatAndCurrentCandle(
			&pastWavetrend, indicators.SpotKlineToMinimalKline(candle),
			10, 21,
		)
		w.waveTrendDat.storeCurrentTci(currentTci)

		once.Do(func() {
			doneC <- nil
		})

		if w.getStopSignal() {
			return
		}
	}
}

func (w *spotWorker) initWaveTrendPastData() error {
	limit := 600
	symbol := w.setting.loadSymbol()
	interval := "1m"

	candles, err := w.bclient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(int(limit)).
		Do(context.Background())
	if err != nil {
		return err
	}
	candles = candles[:len(candles)-1] // drop last candle

	pastWavetrend, _ := indicators.CalculatePastWavetrendData(
		indicators.SpotKlineToMinimalKline(candles),
		10, 21,
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
		10, 21,
	)
	if err != nil {
		return err
	}

	w.waveTrendDat.storePastWaveTrendData(*res)
	return nil
}
