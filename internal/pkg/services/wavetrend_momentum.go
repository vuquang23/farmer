package services

import (
	"context"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/enum"
	"farmer/pkg/errors"
)

type IWavetrendMomentumService interface {
	Calculate(ctx context.Context, market enum.Market, symbolList []string, interval string) ([]*entities.WavetrendMomentum, *errors.DomainError)
}
