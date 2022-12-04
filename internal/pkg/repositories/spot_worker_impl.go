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

	err := r.db.Table("spot_workers").Select("spot_workers.id, spot_workers.symbol", "spot_workers.unit_buy_allowed", "spot_workers.unit_notional", "spot_workers.capital", "SUM(q.unit_bought) AS total_unit_bought").
		Joins("LEFT JOIN (?) q ON q.spot_worker_id = spot_workers.id", query).Group("spot_workers.id").Find(&ret).Error
	if err != nil {
		logger.Error(ctx, err)
		return nil, errors.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}

func (r *spotWorkerRepository) UpdateUnitNotionalByID(ctx context.Context, ID uint64, val float64) *errors.InfraError {
	err := r.db.Table("spot_workers").
		Where("id = ?", ID).
		Update("unit_notional", gorm.Expr("unit_notional + ?", val)).Error
	if err != nil {
		logger.Error(ctx, err)
		return errors.NewInfraErrorDBUpdate(err)
	}

	return nil
}

func (r *spotWorkerRepository) Create(ctx context.Context, w *entities.SpotWorker) (*entities.SpotWorker, *errors.InfraError) {
	err := r.db.Create(w).Error
	if err != nil {
		logger.Error(ctx, err)
		return nil, errors.NewInfraErrorDBInsert(err)
	}
	return w, nil
}

func (r *spotWorkerRepository) DeleteByID(ctx context.Context, ID uint64) *errors.InfraError {
	err := r.db.Delete(&entities.SpotWorker{}, ID).Error
	if err != nil {
		logger.Error(ctx, err)
		return errors.NewInfraErrorDBDelete(err)
	}
	return nil
}

func (r *spotWorkerRepository) AddCapital(ctx context.Context, params *entities.AddCapitalParams) *errors.InfraError {
	logger.Info(ctx, "[AddCapital] update capital in DB")

	err := r.db.Table("spot_workers").
		Where("symbol = ?", params.Symbol).
		Update("capital", gorm.Expr("capital + ?", params.Capital)).
		Update("unit_notional", gorm.Expr("unit_notional + ?/unit_buy_allowed", params.Capital)).Error
	if err != nil {
		logger.Error(ctx, err)
		return errors.NewInfraErrorDBUpdate(err)
	}
	return nil
}
