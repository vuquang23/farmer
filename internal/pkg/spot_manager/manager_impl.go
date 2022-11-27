package spotmanager

import (
	goctx "context"
	"errors"
	"fmt"
	"sync"

	"github.com/adshao/go-binance/v2"

	bn "farmer/internal/pkg/binance"
	"farmer/internal/pkg/constants"
	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/repositories"
	sw "farmer/internal/pkg/spot_worker"
	"farmer/internal/pkg/utils/context"
	"farmer/internal/pkg/utils/logger"
	wtp "farmer/internal/pkg/wavetrend"
	errPkg "farmer/pkg/errors"
)

type spotManager struct {
	mu              *sync.Mutex
	bclient         *binance.Client
	mapSymbolWorker map[string]sw.ISpotWorker // eg for symbol: BTCUSDT, ETHUSDT...
	swRepo          repositories.ISpotWorkerRepository
	mapExchangeInfo map[string]entities.SpotExchangeInfo
}

var manager *spotManager

func InitSpotManager(bclient *binance.Client, swRepo repositories.ISpotWorkerRepository) {
	if manager == nil {
		manager = &spotManager{
			mu:              &sync.Mutex{},
			bclient:         bclient,
			mapSymbolWorker: make(map[string]sw.ISpotWorker),
			swRepo:          swRepo,
			mapExchangeInfo: make(map[string]entities.SpotExchangeInfo),
		}
	}
}

func SpotManagerInstance() ISpotManager {
	return manager
}

func (m *spotManager) Run(ctx goctx.Context, startC chan<- error) {
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

func (m *spotManager) startWorkers(ctx goctx.Context) error {
	workerStatus, err := m.swRepo.GetAllWorkerStatus(ctx)
	if err != nil {
		return err
	}

	for _, w := range workerStatus {
		if err := m.startWorker(ctx, w); err != nil {
			return err
		}
	}

	logger.Infof(ctx, "[startWorkers] start %d workers", len(workerStatus))

	return nil
}

func (m *spotManager) startWorker(ctx goctx.Context, w *entities.SpotWorkerStatus) error {
	worker := sw.NewSpotWorker(
		w.ID,
		bn.BinanceSpotClientInstance(),
		wtp.WavetrendProviderInstance(),
		repositories.SpotTradeRepositoryInstance(),
		repositories.SpotWorkerRepositoryInstance(),
	)
	worker.SetWorkerSettingAndStatus(ctx, *w)
	m.mapSymbolWorker[w.Symbol] = worker

	startC := make(chan error)
	go worker.Run(context.Child(ctx, fmt.Sprintf("[spot-worker] %s", w.Symbol)), startC)
	if err := <-startC; err != nil {
		return err
	}

	return nil
}

func (m *spotManager) CheckHealth() map[string]string {
	ret := make(map[string]string)

	for K, V := range m.mapSymbolWorker {
		ret[K] = V.GetHealth().String()
	}

	return ret
}

func (m *spotManager) CreateNewWorker(ctx goctx.Context, params *entities.CreateNewSpotWorkerParams) error {
	exchangeInfo, ok := m.GetExchangeInfo(params.Symbol)
	if !ok {
		return errPkg.NewDomainErrorNotFound(nil, constants.FieldSymbol)
	}

	notionalOnDownTrend := params.UnitNotional * constants.UnitBuyOnDowntrend
	notionalOnUpTrend := params.UnitNotional * constants.UnitBuyOnUpTrend
	if exchangeInfo.MinNotional > notionalOnDownTrend && exchangeInfo.MinNotional > notionalOnUpTrend {
		err := fmt.Errorf(
			"err constraint on min notional. min notional: %v | notional on downtrend: %v | notional on uptrend: %v",
			exchangeInfo.MinNotional, notionalOnDownTrend, notionalOnUpTrend,
		)
		logger.Error(ctx, err)
		return err
	}

	w, infraErr := m.swRepo.Create(ctx, &entities.SpotWorker{
		Symbol:         params.Symbol,
		UnitBuyAllowed: params.UnitBuyAllowed,
		UnitNotional:   params.UnitNotional,
		Capital:        params.UnitNotional * float64(params.UnitBuyAllowed),
	})
	if infraErr != nil {
		return infraErr
	}

	if err := m.startWorker(ctx, &entities.SpotWorkerStatus{SpotWorker: *w}); err != nil {
		// do not expected error here
		if infraErr := m.swRepo.DeleteByID(ctx, w.ID); infraErr != nil {
			return infraErr
		}

		return err
	}

	return nil
}

func (m *spotManager) StopBot(ctx goctx.Context, params *entities.StopBotParams) error {
	w, ok := m.mapSymbolWorker[params.Symbol]
	if !ok {
		return errors.New("invalid symbol")
	}
	w.SetStopSignal(ctx)
	return nil
}

func (m *spotManager) AddCapital(ctx goctx.Context, params *entities.AddCapitalParams) error {
	w, ok := m.mapSymbolWorker[params.Symbol]
	if !ok {
		return errors.New("invalid symbol")
	}

	infraErr := m.swRepo.AddCapital(ctx, params)
	if infraErr != nil {
		return infraErr
	}

	w.AddCapital(ctx, params.Capital)

	return nil
}
