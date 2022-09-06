package telebot

import (
	"encoding/json"
	"math"

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
	N := 10000.

	for _, e := range en {
		p = append(p, &SpotPairInfo{
			Symbol:                 e.Symbol,
			UsdBenefit:             math.Round(e.UsdBenefit*N) / N,
			BaseAmount:             math.Round(e.BaseAmount*N) / N,
			QuoteAmount:            math.Round(e.QuoteAmount*N) / N,
			CurrentUsdValue:        math.Round(e.CurrentUsdValue*N) / N,
			CurrentUsdValueChanged: math.Round(e.CurrentUsdValueChanged*N) / N,
			UnitBuyAllowed:         e.UnitBuyAllowed,
			UnitNotional:           e.UnitNotional,
			TotalUnitBought:        e.TotalUnitBought,
		})

		totalUsdBenefit += e.UsdBenefit
		currentTotalUsdValueChanged += e.CurrentUsdValueChanged
	}

	return &GetSpotAccountInfoResponse{
		Pairs:                       p,
		TotalUsdBenefit:             math.Round(totalUsdBenefit*N) / N,
		CurrentTotalUsdValueChanged: math.Round(currentTotalUsdValueChanged*N) / N,
	}
}
