package spotmanager

import (
	"context"

	"farmer/internal/pkg/entities"
)

type ISpotManager interface {
	Run(ctx context.Context, startC chan<- error)

	CheckHealth() map[string]string
	CreateNewWorker(ctx context.Context, params *entities.CreateNewSpotWorkerParams) error
	StopWorker(ctx context.Context, params *entities.StopBotParams) error
	AddCapital(ctx context.Context, params *entities.AddCapitalParams) error
}
