package telebot

import (
	"encoding/json"

	"github.com/tidwall/pretty"
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/utils/logger"
)

func (tlb *teleBot) getSpotAccountInfo(ctx tb.Context) {
	log := logger.WithDescription("Get Spot Account Info")
	log.Info("Receive a request")

	var msgResponse string

	ret, err := tlb.spotTradeSvc.GetTradingPairsInfo()
	if err != nil {
		log.Sugar().Error(err)

		byteRes, _ := json.Marshal(err)
		msgResponse = string(pretty.Pretty(byteRes))
		ctx.Send(msgResponse)
		return
	}

	dtoRes := toGetSpotAccountInfoResponse(ret)
	byteRes, _ := json.Marshal(dtoRes)
	msgResponse = string(pretty.Pretty(byteRes))
	ctx.Send(msgResponse)
}

func toGetSpotAccountInfoResponse(en []*entities.SpotTradingPairInfo) *GetSpotAccountInfoResponse {
	p := []*SpotPairInfo{}
	totalUsdBenefit := 0.0
	currentTotalUsdValueChanged := 0.0

	for _, e := range en {
		p = append(p, &SpotPairInfo{
			Symbol:                 e.Symbol,
			UsdBenefit:             e.UsdBenefit,
			BaseAmount:             e.BaseAmount,
			QuoteAmount:            e.QuoteAmount,
			CurrentUsdValue:        e.CurrentUsdValue,
			CurrentUsdValueChanged: e.CurrentUsdValueChanged,
			UnitBuyAllowed:         e.UnitBuyAllowed,
			UnitNotional:           e.UnitNotional,
			TotalUnitBought:        e.TotalUnitBought,
		})

		totalUsdBenefit += e.UsdBenefit
		currentTotalUsdValueChanged += e.CurrentUsdValueChanged
	}

	return &GetSpotAccountInfoResponse{
		Pairs:                       p,
		TotalUsdBenefit:             totalUsdBenefit,
		CurrentTotalUsdValueChanged: currentTotalUsdValueChanged,
	}
}
