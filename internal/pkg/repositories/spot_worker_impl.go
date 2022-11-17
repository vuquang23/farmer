package repositories

import (
	"context"

	"gorm.io/gorm"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/utils/logger"
	"farmer/pkg/errors"
)

type spotWorkerRepository struct {
	db *gorm.DB
}

var spotWorkerRepo *spotWorkerRepository

func InitSpotWorkerRepository(db *gorm.DB) {
	if spotWorkerRepo == nil {
		spotWorkerRepo = &spotWorkerRepository{
			db: db,
		}
	}
}

func SpotWorkerRepositoryInstance() ISpotWorkerRepository {
	return spotWorkerRepo
}

func (r *spotWorkerRepository) GetAllWorkers(ctx context.Context) ([]*entities.SpotWorker, *errors.InfraError) {
	var ret []*entities.SpotWorker
	if err := r.db.Find(&ret).Error; err != nil {
		logger.Error(ctx, err)
		return nil, errors.NewInfraErrorDBSelect(err)
	}
	return ret, nil
}

func (r *spotWorkerRepository) GetAllWorkerStatus(ctx context.Context) ([]*entities.SpotWorkerStatus, *errors.InfraError) {
	var ret []*entities.SpotWorkerStatus

	query := r.db.Table("spot_trades").Select("spot_worker_id, unit_bought").
		Where("side = ? AND is_done = ?", "BUY", false)

	err := r.db.Table("spot_workers").Select("spot_workers.id, spot_workers.symbol", "spot_workers.unit_buy_allowed", "spot_workers.unit_notional", "SUM(q.unit_bought) AS total_unit_bought").
		Joins("LEFT JOIN (?) q ON q.spot_worker_id = spot_workers.id", query).Group("spot_workers.id").Find(&ret).Error
	if err != nil {
		logger.Error(ctx, err)
		return nil, errors.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}

func (r *spotWorkerRepository) UpdateUnitNotionalByID(ctx context.Context, ID uint64, val float64) *errors.InfraError {
	err := r.db.Table("spot_workers").
		Where("id = ?").
		Update("unit_notional", gorm.Expr("unit_notional + ?", val)).Error
	if err != nil {
		logger.Error(ctx, err)
		return errors.NewInfraErrorDBUpdate(err)
	}

	return nil
}
