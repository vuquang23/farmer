package spotworker

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/repositories"
	wt "farmer/internal/pkg/wavetrend"
)

type spotWorker struct {
	ID                  uint64
	bclient             *binance.Client
	exchangeInf         *exchangeInfo
	setting             *workerSetting
	stopSignal          *uint32
	wavetrendProvider   wt.IWavetrendProvider
	wavetrendTimeFrames []string
	stt                 *status
	spotTradeRepo       repositories.ISpotTradeRepository
}

func NewSpotWorker(
	ID uint64, bclient *binance.Client,
	wavetrendProvider wt.IWavetrendProvider,
	spotTradeRepo repositories.ISpotTradeRepository,
) ISpotWorker {
	stopSignal := uint32(0)

	return &spotWorker{
		ID:                  ID,
		bclient:             bclient,
		exchangeInf:         newExchangeInfo(),
		setting:             newWorkerSetting(),
		stopSignal:          &stopSignal,
		wavetrendProvider:   wavetrendProvider,
		wavetrendTimeFrames: []string{"1m", "1h"},
		stt:                 newStatus(),
		spotTradeRepo:       spotTradeRepo,
	}
}

func (w *spotWorker) SetExchangeInfo(info entities.SpotExchangeInfo) error {
	w.exchangeInf.store(info)
	return nil
}

func (w *spotWorker) SetWorkerSettingAndStatus(s entities.SpotWorkerStatus) error {
	w.setting.store(s)
	w.stt.updateTotalUnitBought(int64(s.TotalUnitBought))
	return nil
}

func (w *spotWorker) SetStopSignal() {
	atomic.StoreUint32(w.stopSignal, 1)

	// pass signal to providers
	for _, timeFrame := range w.wavetrendTimeFrames {
		w.wavetrendProvider.SetStopSignal(wavetrendSvcName(w.setting.symbol, timeFrame))
	}
}

func (w *spotWorker) GetHealth() time.Duration {
	return time.Since(w.stt.loadHealth())
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
	return fmt.Sprintf("%s:%s", symbol, timeFrame)
}
