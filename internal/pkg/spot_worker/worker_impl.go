package spotworker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/repositories"
	"farmer/internal/pkg/utils/logger"
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
	spotWorkerRepo      repositories.ISpotWorkerRepository
}

func NewSpotWorker(
	ID uint64, bclient *binance.Client,
	wavetrendProvider wt.IWavetrendProvider,
	spotTradeRepo repositories.ISpotTradeRepository,
	spotWorkerRepo repositories.ISpotWorkerRepository,
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
		spotWorkerRepo:      spotWorkerRepo,
	}
}

func (w *spotWorker) SetExchangeInfo(ctx context.Context, info entities.SpotExchangeInfo) error {
	w.exchangeInf.store(info)
	return nil
}

func (w *spotWorker) SetWorkerSettingAndStatus(ctx context.Context, s entities.SpotWorkerStatus) error {
	w.setting.store(s)
	w.stt.updateTotalUnitBought(int64(s.TotalUnitBought))
	return nil
}

func (w *spotWorker) SetStopSignal(ctx context.Context) {
	atomic.StoreUint32(w.stopSignal, 1)

	// pass signal to wavetrend providers
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

func (w *spotWorker) Run(ctx context.Context, startC chan<- error) {
	for _, timeFrame := range w.wavetrendTimeFrames {
		if err := w.wavetrendProvider.StartService(
			ctx, wavetrendSvcName(w.setting.symbol, timeFrame),
		); err != nil {
			startC <- err
			return
		}
	}

	startC <- nil

	w.runMainProcessor(ctx)
}

func (w *spotWorker) AddCapital(ctx context.Context, capital float64) {
	logger.Info(ctx, "[AddCapital] update capital in memory")
	w.setting.addCapital(capital)
}

func wavetrendSvcName(symbol string, timeFrame string) string {
	return fmt.Sprintf("%s:%s", symbol, timeFrame)
}
