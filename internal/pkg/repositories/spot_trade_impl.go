package repositories

import (
	"time"

	"gorm.io/gorm"

	"farmer/internal/pkg/entities"
	pkgErr "farmer/pkg/errors"
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

func (r *spotTradeRepository) GetNotDoneBuyOrdersByWorkerID(workerID uint64) ([]*entities.SpotTrade, *pkgErr.InfraError) {
	ret := []*entities.SpotTrade{}

	if err := r.db.Where("spot_worker_id = ? AND SIDE = ? AND is_done = ?", workerID, "BUY", false).Find(&ret).Error; err != nil {
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

func (r *spotTradeRepository) UpdateBuyOrders(IDs []uint64) *pkgErr.InfraError {
	if err := r.db.Table("spot_trades").Where("id IN ?", IDs).Update("is_done", true).Error; err != nil {
		return pkgErr.NewInfraErrorDBUpdate(err)
	}

	return nil
}

func (r *spotTradeRepository) CreateSellOrders(spotTrades []*entities.SpotTrade) *pkgErr.InfraError {
	if err := r.db.CreateInBatches(&spotTrades, 100).Error; err != nil {
		return pkgErr.NewInfraErrorDBInsert(err)
	}

	return nil
}

func (r *spotTradeRepository) GetNotDoneBuyOrdersByWorkerIDAndCreatedAt(workerID uint64, createdAfter time.Time) ([]*entities.SpotTrade, *pkgErr.InfraError) {
	ret := []*entities.SpotTrade{}

	err := r.db.
		Where("spot_worker_id = ? AND SIDE = ? AND is_done = ? AND created_at >= ?", workerID, "BUY", false, createdAfter).
		Find(&ret).Error
	if err != nil {
		return nil, pkgErr.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}
