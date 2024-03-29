package services

import (
	"context"

	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
)

type ISpotTradeService interface {
	GetTradingPairsInfo(ctx context.Context) ([]*entities.SpotTradingPairInfo, *pkgErr.DomainError)

	ArchiveTradingData(ctx context.Context, params *entities.ArchiveSpotTradingDataParams) *pkgErr.DomainError
}
