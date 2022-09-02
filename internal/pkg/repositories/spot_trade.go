package repositories

import (
	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
)

type ISpotTradeRepository interface {
	GetNotDoneBuyOrdersBySymbol(symbol string) ([]*entities.SpotTrade, *pkgErr.InfraError)

	CreateBuyOrder(spotTrade entities.SpotTrade) *pkgErr.InfraError
}
