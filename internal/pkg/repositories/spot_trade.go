package repositories

import (
	"time"

	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
)

type ISpotTradeRepository interface {
	GetNotDoneBuyOrdersByWorkerID(workerID uint64) ([]*entities.SpotTrade, *pkgErr.InfraError)
	GetNotDoneBuyOrdersByWorkerIDAndCreatedAtGT(workerID uint64, createdAfter time.Time) ([]*entities.SpotTrade, *pkgErr.InfraError)
	GetTotalQuoteBenefit(workerID uint64) (float64, *pkgErr.InfraError)
	GetBaseAmountAndTotalUnitBought(workerID uint64) (float64, uint64, *pkgErr.InfraError)

	CreateBuyOrder(spotTrade entities.SpotTrade) *pkgErr.InfraError
	CreateSellOrders(spotTrades []*entities.SpotTrade) *pkgErr.InfraError
}
