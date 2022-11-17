package repositories

import (
	"context"
	"time"

	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
)

type ISpotTradeRepository interface {
	GetNotDoneBuyOrdersByWorkerID(ctx context.Context, workerID uint64) ([]*entities.SpotTrade, *pkgErr.InfraError)
	GetNotDoneBuyOrdersByWorkerIDAndCreatedAtGT(workerID uint64, createdAfter time.Time) ([]*entities.SpotTrade, *pkgErr.InfraError)
	GetTotalQuoteBenefit(workerID uint64) (float64, *pkgErr.InfraError)
	GetBaseAmountAndTotalUnitBought(workerID uint64) (float64, uint64, *pkgErr.InfraError)

	CreateBuyOrder(ctx context.Context, spotTrade entities.SpotTrade) *pkgErr.InfraError
	CreateSellOrders(ctx context.Context, spotTrades []*entities.SpotTrade) *pkgErr.InfraError
}
