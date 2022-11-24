package services

import (
	"context"

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

func (s *spotTradeService) GetTradingPairsInfo(ctx context.Context) ([]*en.SpotTradingPairInfo, *pkgErr.DomainError) {
	var ret []*en.SpotTradingPairInfo

	workers, infraErr := s.spotWorkerRepo.GetAllWorkers(ctx)
	if infraErr != nil {
		return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	for _, w := range workers {
		info := &en.SpotTradingPairInfo{
			Symbol:         w.Symbol,
			UnitBuyAllowed: w.UnitBuyAllowed,
			UnitNotional:   w.UnitNotional,
		}

		usdBenefit, infraErr := s.spotTradeRepo.GetTotalQuoteBenefit(w.ID)
		if infraErr != nil {
			return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
		}
		info.UsdBenefit = usdBenefit

		baseAmount, totalUnitBought, infraErr := s.spotTradeRepo.GetBaseAmountAndTotalUnitBought(w.ID)
		if infraErr != nil {
			return nil, pkgErr.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
		}
		info.BaseAmount = baseAmount
		info.TotalUnitBought = totalUnitBought

		price, domainErr := b.GetSpotPrice(ctx, s.bclient, w.Symbol)
		if domainErr != nil {
			return nil, domainErr
		}
		info.QuoteAmount = info.UnitNotional * (float64(info.UnitBuyAllowed) - float64(info.TotalUnitBought))
		info.CurrentUsdValue = info.QuoteAmount + info.BaseAmount*price

		ret = append(ret, info)
	}

	return ret, nil
}
