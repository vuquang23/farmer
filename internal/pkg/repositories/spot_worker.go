package repositories

import (
	"farmer/internal/pkg/entities"
	"farmer/pkg/errors"
)

type ISpotWorkerRepository interface {
	GetAllWorkers() ([]*entities.SpotWorker, *errors.InfraError)
	GetAllWorkerStatus() ([]*entities.SpotWorkerStatus, *errors.InfraError)
}
