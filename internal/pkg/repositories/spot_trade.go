package repositories

import (
	"time"

	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
)

type ISpotTradeRepository interface {
	GetNotDoneBuyOrdersByWorkerID(workerID uint64) ([]*entities.SpotTrade, *pkgErr.InfraError)
	GetNotDoneBuyOrdersByWorkerIDAndCreatedAt(workerID uint64, createdAfter time.Time) ([]*entities.SpotTrade, *pkgErr.InfraError)

	CreateBuyOrder(spotTrade entities.SpotTrade) *pkgErr.InfraError
	CreateSellOrders(spotTrades []*entities.SpotTrade) *pkgErr.InfraError

	UpdateBuyOrders(IDs []uint64) *pkgErr.InfraError
}
