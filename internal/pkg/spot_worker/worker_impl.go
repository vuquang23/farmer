package spotworker

import (
	"fmt"
	"sync/atomic"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/entities"
	wt "farmer/internal/pkg/wavetrend"
)

type spotWorker struct {
	bclient             *binance.Client
	exchangeInf         *exchangeInfo
	setting             *workerSetting
	stopSignal          *uint32
	wavetrendProvider   wt.IWavetrendProvider
	wavetrendTimeFrames []string
}

func NewSpotWorker(bclient *binance.Client, wavetrendProvider wt.IWavetrendProvider) ISpotWorker {
	stopSignal := uint32(0)

	return &spotWorker{
		bclient:             bclient,
		exchangeInf:         newExchangeInfo(),
		setting:             newWorkerSetting(),
		stopSignal:          &stopSignal,
		wavetrendProvider:   wavetrendProvider,
		wavetrendTimeFrames: []string{"1m", "1h"},
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
	for _, timeFrame := range w.wavetrendTimeFrames {
		if err := w.wavetrendProvider.StartService(
			wavetrendSvcName(w.setting.symbol, timeFrame),
		); err != nil {
			startC <- err
			return
		}
	}

	startC <- nil

	w.runMainProcessor()
}

func wavetrendSvcName(symbol string, timeFrame string) string {
	return fmt.Sprintf("%s-%s", symbol, timeFrame)
}
