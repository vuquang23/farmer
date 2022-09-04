package telebot

import (
	"encoding/json"

	"github.com/tidwall/pretty"
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/entities"
)

func (tlb *teleBot) getSpotAccountInfo(ctx tb.Context) {
	var msgResponse string

	ret, err := tlb.spotTradeSvc.GetTradingPairsInfo()
	if err != nil {
		byteRes, _ := json.Marshal(err)
		msgResponse = string(pretty.Pretty(byteRes))
		ctx.Send(msgResponse)
		return
	}

	dtoRes := toGetAccountInfoResponse(ret)
	byteRes, _ := json.Marshal(dtoRes)
	msgResponse = string(pretty.Pretty(byteRes))
	ctx.Send(msgResponse)
}

func toGetAccountInfoResponse(en []*entities.TradingPairInfo) *GetAccountInfoResponse {
	p := []*PairInfo{}
	totalUsdBenefit := 0.0
	currentTotalUsdValueChanged := 0.0

	for _, e := range en {
		p = append(p, &PairInfo{
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

	return &GetAccountInfoResponse{
		Pairs:                       p,
		TotalUsdBenefit:             totalUsdBenefit,
		CurrentTotalUsdValueChanged: currentTotalUsdValueChanged,
	}
}
