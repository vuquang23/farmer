package spotmanager

import sb "farmer/internal/pkg/spot_binance"

type SpotManager struct {
	mapSymbolWorker map[string]*sb.SpotWorker
}

var spotManager *SpotManager

func InitSpotManager() {
	if spotManager == nil {
		spotManager = &SpotManager{}
	}
}

func SpotManagerInstance() ISpotManager {
	return spotManager
}
