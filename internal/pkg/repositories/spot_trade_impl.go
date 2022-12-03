package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/utils/logger"
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

func (r *spotTradeRepository) GetNotDoneBuyOrdersByWorkerID(ctx context.Context, workerID uint64) ([]*entities.SpotTrade, *pkgErr.InfraError) {
	ret := []*entities.SpotTrade{}

	if err := r.db.Where("spot_worker_id = ? AND side = ? AND is_done = ?", workerID, "BUY", false).Find(&ret).Error; err != nil {
		logger.Error(ctx, err)
		return nil, pkgErr.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}

func (r *spotTradeRepository) CreateBuyOrder(ctx context.Context, spotTrade entities.SpotTrade) *pkgErr.InfraError {
	if err := r.db.Create(&spotTrade).Error; err != nil {
		logger.Error(ctx, err)
		return pkgErr.NewInfraErrorDBInsert(err)
	}

	return nil
}

func (r *spotTradeRepository) CreateSellOrders(ctx context.Context, spotTrades []*entities.SpotTrade) *pkgErr.InfraError {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&spotTrades, 100).Error; err != nil {
			logger.Error(ctx, err)
			return pkgErr.NewInfraErrorDBInsert(err)
		}

		buyOrderIDs := make([]uint64, len(spotTrades))
		for idx, t := range spotTrades {
			buyOrderIDs[idx] = t.Ref
		}

		if err := tx.Table("spot_trades").Where("id IN ?", buyOrderIDs).Update("is_done", true).Error; err != nil {
			logger.Error(ctx, err)
			return pkgErr.NewInfraErrorDBUpdate(err)
		}

		return nil
	})

	if err != nil {
		infraErr, ok := err.(*pkgErr.InfraError)
		if ok {
			return infraErr
		}
		return pkgErr.NewInfraErrorDBUnknown(err)
	}

	return nil
}

func (r *spotTradeRepository) GetNotDoneBuyOrdersByWorkerIDAndCreatedAtGT(workerID uint64, createdAfter time.Time) ([]*entities.SpotTrade, *pkgErr.InfraError) {
	ret := []*entities.SpotTrade{}

	err := r.db.
		Where("spot_worker_id = ? AND side = ? AND is_done = ? AND created_at >= ?", workerID, "BUY", false, createdAfter).
		Find(&ret).Error
	if err != nil {
		return nil, pkgErr.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}

func (r *spotTradeRepository) GetTotalQuoteBenefit(workerID uint64) (float64, *pkgErr.InfraError) {
	type response struct {
		TotalQuoteBenefit float64
	}
	var ret response

	querySell := r.db.Table("spot_trades").Where("spot_worker_id = ? AND side = ?", workerID, "SELL")
	err := r.db.Table("spot_trades").Joins("JOIN (?) querySell ON querySell.ref = spot_trades.id", querySell).Group("spot_trades.spot_worker_id").
		Select("SUM(querySell.cummulative_quote_qty - spot_trades.cummulative_quote_qty) as total_quote_benefit").Find(&ret).Error
	if err != nil {
		return 0, pkgErr.NewInfraErrorDBSelect(err)
	}

	return ret.TotalQuoteBenefit, nil
}

func (r *spotTradeRepository) GetAggregatedNotSoldBuyOrders(ctx context.Context, workerID uint64) (*entities.AggregatedBuyOrders, *pkgErr.InfraError) {
	var ret *entities.AggregatedBuyOrders

	err := r.db.Table("spot_trades").Where("spot_worker_id = ? AND side = ? AND is_done = ?", workerID, "BUY", false).
		Group("spot_worker_id").
		Select(`
			SUM(qty) as total_base_amount, 
			SUM(cummulative_quote_qty) as total_cummulative_quote_qty, 
			SUM(unit_bought) as total_unit_bought
		`).
		Find(&ret).Error
	if err != nil {
		return nil, pkgErr.NewInfraErrorDBSelect(err)
	}

	return ret, nil
}

func (r *spotTradeRepository) ArchiveTradingData(ctx context.Context, symbol string) *pkgErr.InfraError {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var spotTrades []entities.SpotTrade
		err := tx.Where("symbol = ? AND is_done = ?", symbol, true).Find(&spotTrades).Error
		if err != nil {
			logger.Error(ctx, err)
			return pkgErr.NewInfraErrorDBSelect(err)
		}

		historySpotTrades := make([]entities.HistorySpotTrade, len(spotTrades))
		for idx, t := range spotTrades {
			historySpotTrades[idx] = entities.NewHistorySpotTrade(t)
		}
		err = tx.Create(historySpotTrades).Error
		if err != nil {
			logger.Error(ctx, err)
			return pkgErr.NewInfraErrorDBInsert(err)
		}

		err = tx.Where("symbol = ?", symbol).Delete(&entities.SpotTrade{}).Error
		if err != nil {
			logger.Error(ctx, err)
			return pkgErr.NewInfraErrorDBDelete(err)
		}

		err = tx.Where("symbol = ?", symbol).Delete(&entities.SpotWorker{}).Error
		if err != nil {
			logger.Error(ctx, err)
			return pkgErr.NewInfraErrorDBDelete(err)
		}

		return nil
	})

	if err != nil {
		infraErr, ok := err.(*pkgErr.InfraError)
		if ok {
			return infraErr
		}
		return pkgErr.NewInfraErrorDBUnknown(err)
	}

	return nil
}
