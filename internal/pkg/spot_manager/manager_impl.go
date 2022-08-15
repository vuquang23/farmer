package spotmanager

import (
	sw "farmer/internal/pkg/spot_worker"

	"github.com/adshao/go-binance/v2"
)

type spotManager struct {
	bclient         *binance.Client
	mapSymbolWorker map[string]sw.ISpotWorker // eg for symbol: BTCUSDT, ETHUSDT...
}

var manager *spotManager

func InitSpotManager(bclient *binance.Client) {
	if manager == nil {
		manager = &spotManager{
			bclient:         bclient,
			mapSymbolWorker: make(map[string]sw.ISpotWorker),
		}
	}
}

func SpotManagerInstance() ISpotManager {
	return manager
}

func (m *spotManager) Run() error {
	go m.updateExchangeInfoPeriodically()

	return nil
}
