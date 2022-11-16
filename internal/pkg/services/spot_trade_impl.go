package services

import (
	"github.com/adshao/go-binance/v2"

	b "farmer/internal/pkg/binance"
	en "farmer/internal/pkg/entities"
	"farmer/internal/pkg/repositories"
	pkgErr "farmer/pkg/errors"
)

type spotTradeService struct {
	bclient        *binance.Client
	spotTradeRepo  repositories.ISpotTradeRepository
	spotWorkerRepo repositories.ISpotWorkerRepository
}

var spotTradeSvc *spotTradeService

func InitSpotTradeService(
	bclient *binance.Client,
	spotTradeRepo repositories.ISpotTradeRepository,
	spotWorkerRepo repositories.ISpotWorkerRepository,
) {
	if spotTradeSvc == nil {
		spotTradeSvc = &spotTradeService{
			bclient:        bclient,
			spotTradeRepo:  spotTradeRepo,
			spotWorkerRepo: spotWorkerRepo,
		}
	}
}

func SpotTradeServiceInstance() ISpotTradeService {
	return spotTradeSvc
}

func (s *spotTradeService) GetTradingPairsInfo() ([]*en.SpotTradingPairInfo, *pkgErr.DomainError) {
	var ret []*en.SpotTradingPairInfo

	workers, err := s.spotWorkerRepo.GetAllWorkers()
	if err != nil {
		return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(err)
	}

	for _, w := range workers {
		info := &en.SpotTradingPairInfo{
			Symbol:         w.Symbol,
			UnitBuyAllowed: w.UnitBuyAllowed,
			UnitNotional:   w.UnitNotional,
		}

		if usdBenefit, err := s.spotTradeRepo.GetTotalQuoteBenefit(w.ID); err != nil {
			return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(err)
		} else {
			info.UsdBenefit = usdBenefit
		}

		if baseAmount, totalUnitBought, err := s.spotTradeRepo.GetBaseAmountAndTotalUnitBought(w.ID); err != nil {
			return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(err)
		} else {
			info.BaseAmount = baseAmount
			info.TotalUnitBought = totalUnitBought
		}

		if price, err := b.GetSpotPrice(s.bclient, w.Symbol); err != nil {
			return nil, err
		} else {
			info.QuoteAmount = info.UnitNotional * (float64(info.UnitBuyAllowed) - float64(info.TotalUnitBought))
			info.CurrentUsdValue = info.QuoteAmount + info.BaseAmount*price + info.UsdBenefit
			info.CurrentUsdValueChanged = info.CurrentUsdValue - info.UnitNotional*float64(info.UnitBuyAllowed)
		}

		ret = append(ret, info)
	}

	return ret, nil
}
