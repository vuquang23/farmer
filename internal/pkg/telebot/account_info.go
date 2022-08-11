package telebot

import (
	"encoding/json"

	tb "gopkg.in/telebot.v3"
)

type AccountInfoQuery struct {
	Exchange string   `json:"exchange"`
	Symbols  []string `json:"symbols"` // eg: BTCUSDT
}

type SpotSymbolDetails struct {
	Symbol       string  `json:"symbol"`
	Benefit      float64 `json:"benefit"`
	SoldOrder    uint64  `json:"soldOrder"`
	NotSoldOrder uint64  `json:"notSoldOrder"`
}

type SpotAccountInfoResponse struct {
	Exchange          string              `json:"exchange"`
	TotalBenefit      float64             `json:"totalBenefit"`      // in USD
	TotalSoldOrder    uint64              `json:"totalSoldOrder"`    // number of orders that is sold successfully.
	TotalNotSoldOrder uint64              `json:"totalNotSoldOrder"` // number of orders that is waiting for good price to sell. not pending
	Details           []SpotSymbolDetails `json:"details"`
}

func getSpotAccountInfo(ctx tb.Context, dto *AccountInfoQuery) {
	fake := &SpotAccountInfoResponse{
		Exchange:          "binance",
		TotalBenefit:      10,
		TotalSoldOrder:    20,
		TotalNotSoldOrder: 30,
		Details: []SpotSymbolDetails{
			{
				Symbol:       "BTCUSDT",
				Benefit:      2,
				SoldOrder:    5,
				NotSoldOrder: 10,
			},
		},
	}

	data, _ := json.Marshal(fake)
	ctx.Send(string(data))
}
