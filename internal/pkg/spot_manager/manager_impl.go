package spotmanager

import sw "farmer/internal/pkg/spot_worker"

type SpotManager struct {
	mapSymbolWorker map[string]sw.ISpotWorker
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
