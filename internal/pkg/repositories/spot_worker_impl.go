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
