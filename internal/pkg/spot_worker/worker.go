package spotworker

import (
	"time"

	"farmer/internal/pkg/entities"
)

type ISpotWorker interface {
	SetExchangeInfo(info entities.SpotExchangeInfo) error
	SetWorkerSettingAndStatus(s entities.SpotWorkerStatus) error
	SetStopSignal()

	//GetHealth return time duration from last update until now.
	GetHealth() time.Duration

	Run(startC chan<- error)
}
