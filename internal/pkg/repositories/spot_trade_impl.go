package repositories

import (
	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"

	"gorm.io/gorm"
)

type spotTradeRepository struct {
	db *gorm.DB
}

var spotTradeRepo *spotTradeRepository

func InitSpotTradeRepository(db *gorm.DB) {
	if spotTradeRepo == nil {
		spotTradeRepo = &spotTradeRepository{
			db: db,
		}
	}
}

func SpotTradeRepositoryInstance() ISpotTradeRepository {
	return spotTradeRepo
}

func (r *spotTradeRepository) GetNotDoneBuyOrdersBySymbol(symbol string) ([]*entities.SpotTrade, *pkgErr.InfraError) {
	ret := []*entities.SpotTrade{}

	if err := r.db.Where("symbol = ? AND is_done = ?", symbol, false).Find(&ret).Error; err != nil {
		return nil, pkgErr.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}

func (r *spotTradeRepository) CreateBuyOrder(spotTrade entities.SpotTrade) *pkgErr.InfraError {
	if err := r.db.Create(&spotTrade).Error; err != nil {
		return pkgErr.NewInfraErrorDBInsert(err)
	}

	return nil
}
