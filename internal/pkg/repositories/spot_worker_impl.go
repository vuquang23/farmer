package repositories

import (
	"gorm.io/gorm"

	"farmer/internal/pkg/entities"
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

func (r *spotWorkerRepository) GetAllWorkers() ([]*entities.SpotWorker, *errors.InfraError) {
	ret := []*entities.SpotWorker{}
	if err := r.db.Find(&ret).Error; err != nil {
		return nil, errors.NewInfraErrorDBSelect(err)
	}
	return ret, nil
}

func (r *spotWorkerRepository) GetAllWorkerStatus() ([]*entities.SpotWorkerStatus, *errors.InfraError) {
	ret := []*entities.SpotWorkerStatus{}

	query := r.db.Table("spot_trades").Select("spot_worker_id,unit_bought").
		Where("side = ? AND is_done = ?", "BUY", false)

	err := r.db.Table("spot_workers").Select("spot_workers.symbol", "spot_workers.unit_buy_allowed", "spot_workers.unit_notional", "SUM(q.unit_bought) AS total_unit_bought").
		Joins("JOIN (?) q ON q.spot_worker_id = spot_workers.id", query).Group("spot_workers.id").Find(&ret).Error
	if err != nil {
		return nil, errors.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}
