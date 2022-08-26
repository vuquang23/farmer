package worker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adshao/go-binance/v2"

	c "farmer/internal/pkg/constants"
	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/indicators"
	"farmer/internal/pkg/utils/logger"
)

type worker struct {
	bclient     *binance.Client
	mu          *sync.Mutex
	symbol      string
	timeFrame   string
	pastData    *pastWavetrendData
	currentData *currentWavetrendData
	setting     *workerSetting
	stopSignal  *uint32
}

//NewWavetrendWorker calculates wavetrend base on svcName (eg: BTCUSDT-1h, ETHUSDT-1m)
func NewWavetrendWorker(svcName string, bclient *binance.Client) IWavetrendWorker {
	s := strings.Split(svcName, "-")
	stopSignal := uint32(0)

	return &worker{
		bclient:     bclient,
		mu:          &sync.Mutex{},
		symbol:      s[0],
		timeFrame:   s[1],
		pastData:    newPastWavetrendData(),
		currentData: newCurrentWavetrendData(),
		setting:     newWorkerSetting(s[1]),
		stopSignal:  &stopSignal,
	}
}

func (w *worker) SetStopSignal() {
	atomic.StoreUint32(w.stopSignal, 1)
}

func (w *worker) getStopSignal() bool {
	return atomic.LoadUint32(w.stopSignal) > 0
}

func (w *worker) GetCurrentTci() float64 {
	return w.loadCurrentTci()
}

func (w *worker) GetCurrentDifWavetrend() float64 {
	return w.loadCurrentDifWavetrend()
}

func (w *worker) Run(done chan<- error) {
	log := logger.WithDescription(fmt.Sprintf("[%s-%s] Update wave trend periodically", w.symbol, w.timeFrame))

	if err := w.initWaveTrendPastData(); err != nil {
		done <- errors.NewDomainErrorInitWavetrendData(err, w.symbol)
		return
	}

	once := &sync.Once{}
	ticker := time.NewTicker(w.setting.sleepAfterQuery)
	periodMilis := w.setting.timeFrameUnixMili
	for ; !w.getStopSignal(); <-ticker.C {
		// check whether now is new interval
		lastOpenTime := w.loadLastOpenTime()
		now := uint64(time.Now().UnixMilli())
		if now-lastOpenTime > periodMilis*2 {
			err := w.updateWaveTrendForNextInterval(
				lastOpenTime+periodMilis,
				(now-lastOpenTime)/periodMilis-1,
			)
			if err != nil {
				log.Sugar().Error(err)
				continue
			}
		}

		// get candle of current interval
		candle, err := w.bclient.NewKlinesService().
			Symbol(w.symbol).
			Interval(w.timeFrame).
			Limit(1).
			Do(context.Background())
		if err != nil {
			log.Sugar().Error(err)
			continue
		}

		pastWavetrend := w.loadPastWaveTrendData()
		currentTci, currentDifWavetrend := indicators.CalcCurrentTciAndDifWavetrend(
			&pastWavetrend, indicators.SpotKlineToMinimalKline(candle),
			c.EmaLenN1, c.EmaLenN2, c.AvgPeriodLen,
		)
		w.storeCurrentTci(currentTci)
		w.storeCurrentDifWavetrend(currentDifWavetrend)

		once.Do(func() {
			done <- nil
		})
	}
}

func (w *worker) updateWaveTrendForNextInterval(fromOpenTime uint64, limit uint64) error {
	symbol := w.symbol
	interval := w.timeFrame

	candles, err := w.bclient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(int64(fromOpenTime)).
		Limit(int(limit)).
		Do(context.Background())
	if err != nil {
		return err
	}

	pastWaveTrend := w.loadPastWaveTrendData()
	res, err := indicators.UpdatePastWavetrendDataWithNewCandles(
		&pastWaveTrend,
		indicators.SpotKlineToMinimalKline(candles),
		c.EmaLenN1, c.EmaLenN2, c.AvgPeriodLen, c.DifWavetrendLen,
	)
	if err != nil {
		return err
	}

	w.storePastWaveTrendData(*res)
	return nil
}

func (w *worker) initWaveTrendPastData() error {
	interval := w.timeFrame

	candles, err := w.bclient.NewKlinesService().
		Symbol(w.symbol).
		Interval(interval).
		Limit(int(c.KlineHistoryLen)).
		Do(context.Background())
	if err != nil {
		return err
	}
	candles = candles[:len(candles)-1] // drop last candle

	pastWavetrend, _ := indicators.CalcPastWavetrendData(
		indicators.SpotKlineToMinimalKline(candles),
		c.EmaLenN1, c.EmaLenN2, w.setting.tciLen, c.AvgPeriodLen, c.DifWavetrendLen,
	)

	w.storePastWaveTrendData(*pastWavetrend)
	return nil
}
