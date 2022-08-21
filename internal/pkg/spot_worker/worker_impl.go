package spotworker

import (
	"sync/atomic"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/entities"
)

type spotWorker struct {
	bclient               *binance.Client
	exchangeInf           *exchangeInfo
	setting               *workerSetting
	waveTrendDat          *waveTrendData
	secondaryWavetrendDat *secondaryWavetrendData

	stopSignal *uint32
}

func NewSpotWorker(bclient *binance.Client) ISpotWorker {
	stopSignal := uint32(0)

	return &spotWorker{
		bclient:               bclient,
		exchangeInf:           newExchangeInfo(),
		setting:               newWorkerSetting(),
		waveTrendDat:          newWaveTrendData(),
		secondaryWavetrendDat: newSecondaryWaveTrendData(),

		stopSignal: &stopSignal,
	}
}

func (w *spotWorker) SetExchangeInfo(info entities.ExchangeInfo) error {
	w.exchangeInf.store(info)
	return nil
}

func (w *spotWorker) SetWorkerSetting(setting entities.SpotWorker) error {
	w.setting.store(setting)
	return nil
}

func (w *spotWorker) SetStopSignal() {
	atomic.StoreUint32(w.stopSignal, 1)
}

func (w *spotWorker) getStopSignal() bool {
	return atomic.LoadUint32(w.stopSignal) > 0
}

func (w *spotWorker) Run(startC chan<- error) {
	doneC := make(chan error)
	go w.updateWaveTrendPeriodically(doneC)
	if err := <-doneC; err != nil {
		startC <- err
		return
	}

	go w.updateSecondaryWavetrendPeriodically(doneC)
	if err := <-doneC; err != nil {
		startC <- err
		return
	}
	startC <- nil

	w.runMainProcessor()
}
