package telebot

import (
	"farmer/internal/pkg/enum"

	tb "gopkg.in/telebot.v3"
)

type CreateSpotBotRequest struct {
	Exchange string  `json:"exchange"`
	Amount   float64 `json:"amount"` // in USDT
}

func createSpotBot(ctx tb.Context, req *CreateSpotBotRequest) {
	ctx.Send("ok create spot bot")
}

type StopBotBotRequest struct {
	Exchange string            `json:"exchange"`
	StopType enum.SpotStopType `json:"stopType"`
}

func stopSpotBot(ctx tb.Context, req *StopBotBotRequest) {
	ctx.Send("ok stop spot bot")
}
