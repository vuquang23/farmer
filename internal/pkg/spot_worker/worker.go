package spotworker

import (
	"context"
	"time"

	"farmer/internal/pkg/entities"
)

type ISpotWorker interface {
	SetExchangeInfo(ctx context.Context, info entities.SpotExchangeInfo) error
	SetWorkerSettingAndStatus(ctx context.Context, s entities.SpotWorkerStatus) error
	SetStopSignal(ctx context.Context)
	AddCapital(ctx context.Context, capital float64)

	//GetHealth return time duration from last update until now.
	GetHealth() time.Duration

	Run(ctx context.Context, startC chan<- error)
}
