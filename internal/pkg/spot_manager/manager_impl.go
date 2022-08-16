package spotmanager

import (
	"github.com/adshao/go-binance/v2"

	bn "farmer/internal/pkg/binance"
	"farmer/internal/pkg/repositories"
	sw "farmer/internal/pkg/spot_worker"
	"farmer/internal/pkg/utils/logger"
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

func (m *spotManager) Run(startC chan<- error) {
	if err := m.loadWorkers(); err != nil {
		startC <- err
		return
	}

	doneC := make(chan error)
	go m.updateExchangeInfoPeriodically(doneC)
	if err := <-doneC; err != nil {
		startC <- err
		return
	}

	if err := m.startWorkers(); err != nil {
		startC <- err
		return
	}

	logger.Logger.Info("Start worker manager successfully")

	startC <- nil
}

func (m *spotManager) loadWorkers() error {
	logger.Logger.Debug("Load workers")

	workerEntities, err := m.swRepo.GetAllWorkers()
	if err != nil {
		return err
	}

	for _, workerEntity := range workerEntities {
		worker := sw.NewSpotWorker(bn.BinanceSpotClientInstance())
		worker.SetWorkerSetting(*workerEntity)
		m.mapSymbolWorker[workerEntity.Symbol] = worker
	}

	return nil
}

func (m *spotManager) startWorkers() error {
	logger.Logger.Debug("Start workers")

	workerEntities, err := m.swRepo.GetAllWorkers()
	if err != nil {
		return err
	}

	// all worker should start OK
	for _, workerEntity := range workerEntities {
		worker := m.mapSymbolWorker[workerEntity.Symbol]
		startC := make(chan error)
		go worker.Run(startC)
		if err := <-startC; err != nil {
			return err
		}
	}

	return nil
}
