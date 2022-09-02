package spotworker

import "farmer/internal/pkg/entities"

type ISpotWorker interface {
	SetExchangeInfo(info entities.ExchangeInfo) error
	SetWorkerSettingAndStatus(s entities.SpotWorkerStatus) error
	SetStopSignal()

	Run(startC chan<- error)
}
