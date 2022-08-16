package telebot

import (
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/enum"
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
