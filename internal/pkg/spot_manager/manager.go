package spotmanager

import (
	"context"

	"farmer/internal/pkg/entities"
)

type ISpotManager interface {
	Run(ctx context.Context, startC chan<- error)

	CheckHealth() map[string]string
	CreateNewWorker(ctx context.Context, params *entities.CreateNewSpotWorkerParams) error
	StopWorker(ctx context.Context, params *entities.StopWorkerParams) error
	AddCapital(ctx context.Context, params *entities.AddCapitalParams) error

	IsActiveWorker(ctx context.Context, symbol string) bool
}
