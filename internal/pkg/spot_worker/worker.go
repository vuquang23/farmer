package spotworker

import "farmer/internal/pkg/entities"

type ISpotWorker interface {
	SetExchangeInfo(info entities.SpotExchangeInfo) error
	SetWorkerSettingAndStatus(s entities.SpotWorkerStatus) error
	SetStopSignal()

	Run(startC chan<- error)
}
