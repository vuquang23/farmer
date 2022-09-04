package services

import (
	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
)

type ISpotTradeService interface {
	GetTradingPairsInfo() ([]*entities.TradingPairInfo, *pkgErr.DomainError)
}
