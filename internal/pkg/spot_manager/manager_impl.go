package spotmanager

import (
	"github.com/adshao/go-binance/v2"

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

	doneC := make(chan struct{})
	go m.updateExchangeInfoPeriodically(doneC)
	<-doneC

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
		worker := sw.NewSpotWorker()
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
		if err := worker.Run(); err != nil {
			return err
		}
	}

	return nil
}
