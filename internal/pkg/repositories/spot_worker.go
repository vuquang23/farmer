package repositories

import (
	"context"

	"farmer/internal/pkg/entities"
	"farmer/pkg/errors"
)

type ISpotWorkerRepository interface {
	GetAllWorkers(ctx context.Context) ([]*entities.SpotWorker, *errors.InfraError)
	GetAllWorkerStatus(ctx context.Context) ([]*entities.SpotWorkerStatus, *errors.InfraError)

	UpdateUnitNotionalByID(ctx context.Context, ID uint64, val float64) *errors.InfraError
}
