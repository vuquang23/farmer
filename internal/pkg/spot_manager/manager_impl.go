package spotmanager

import (
	goctx "context"
	"fmt"

	"github.com/adshao/go-binance/v2"

	bn "farmer/internal/pkg/binance"
	"farmer/internal/pkg/repositories"
	sw "farmer/internal/pkg/spot_worker"
	"farmer/internal/pkg/utils/context"
	"farmer/internal/pkg/utils/logger"
	wtp "farmer/internal/pkg/wavetrend"
)

type spotManager struct {
	bclient         *binance.Client
	mapSymbolWorker map[string]sw.ISpotWorker // eg for symbol: BTCUSDT, ETHUSDT...
	swRepo          repositories.ISpotWorkerRepository
}

var manager *spotManager

func InitSpotManager(bclient *binance.Client, swRepo repositories.ISpotWorkerRepository) {
	if manager == nil {
		manager = &spotManager{
			bclient:         bclient,
			mapSymbolWorker: make(map[string]sw.ISpotWorker),
			swRepo:          swRepo,
		}
	}
}

func SpotManagerInstance() ISpotManager {
	return manager
}

func (m *spotManager) Run(ctx goctx.Context, startC chan<- error) {
	if err := m.loadWorkers(ctx); err != nil {
		startC <- err
		return
	}

	doneC := make(chan error)
	go m.updateExchangeInfoPeriodically(context.Child(ctx, "spot manager update exchange info periodically"), doneC)
	if err := <-doneC; err != nil {
		startC <- err
		return
	}

	if err := m.startWorkers(ctx); err != nil {
		startC <- err
		return
	}

	logger.Info(ctx, "[Run] start worker manager successfully")

	startC <- nil
}

func (m *spotManager) loadWorkers(ctx goctx.Context) error {
	logger.Info(ctx, "[loadWorkers] start to load workers")

	workerStatus, err := m.swRepo.GetAllWorkerStatus(ctx)
	if err != nil {
		return err
	}

	for _, w := range workerStatus {
		worker := sw.NewSpotWorker(
			w.ID,
			bn.BinanceSpotClientInstance(),
			wtp.WavetrendProviderInstance(),
			repositories.SpotTradeRepositoryInstance(),
			repositories.SpotWorkerRepositoryInstance(),
		)
		worker.SetWorkerSettingAndStatus(*w)
		m.mapSymbolWorker[w.Symbol] = worker
	}

	return nil
}

func (m *spotManager) startWorkers(ctx goctx.Context) error {
	workerEntities, err := m.swRepo.GetAllWorkers(ctx)
	if err != nil {
		return err
	}

	// all worker should start OK
	for _, workerEntity := range workerEntities {
		worker := m.mapSymbolWorker[workerEntity.Symbol]
		startC := make(chan error)
		go worker.Run(context.Child(ctx, fmt.Sprintf("[spot-worker] %s", workerEntity.Symbol)), startC)
		if err := <-startC; err != nil {
			return err
		}
	}

	logger.Infof(ctx, "[startWorkers] start %d workers", len(workerEntities))

	return nil
}

func (m *spotManager) CheckHealth() map[string]string {
	ret := make(map[string]string)

	for K, V := range m.mapSymbolWorker {
		ret[K] = V.GetHealth().String()
	}

	return ret
}
