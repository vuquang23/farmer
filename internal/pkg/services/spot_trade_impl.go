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
	ret := []*en.SpotTradingPairInfo{}

	workers, err := s.spotWorkerRepo.GetAllWorkers()
	if err != nil {
		return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(err)
	}

	for _, w := range workers {
		temp := &en.SpotTradingPairInfo{
			Symbol:         w.Symbol,
			UnitBuyAllowed: w.UnitBuyAllowed,
			UnitNotional:   w.UnitNotional,
		}

		if usdBenefit, err := s.spotTradeRepo.GetTotalQuoteBenefit(w.ID); err != nil {
			return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(err)
		} else {
			temp.UsdBenefit = usdBenefit
		}

		if baseAmount, totalUnitBought, err := s.spotTradeRepo.GetBaseAmountAndTotalUnitBought(w.ID); err != nil {
			return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(err)
		} else {
			temp.BaseAmount = baseAmount
			temp.TotalUnitBought = totalUnitBought
		}

		if price, err := b.GetSpotPrice(s.bclient, w.Symbol); err != nil {
			return nil, err
		} else {
			temp.QuoteAmount = temp.UnitNotional * (float64(temp.UnitBuyAllowed) - float64(temp.TotalUnitBought))
			temp.CurrentUsdValue = temp.QuoteAmount + temp.BaseAmount*price
			temp.CurrentUsdValueChanged = temp.CurrentUsdValue - temp.UnitNotional*float64(temp.UnitBuyAllowed)
		}

		ret = append(ret, temp)
	}

	return ret, nil
}
