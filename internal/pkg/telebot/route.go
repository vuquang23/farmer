package telebot

import (
	"strings"

	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/config"
)

const (
	GetSpotAccountInfoCmd = "get!/spot/account-info"
	GetSpotHealthCmd      = "get!/spot/health"

	CreateSpotWorkerCmd       = "post!/spot"
	AddCapitalSpotWorkerCmd   = "post!/spot/add-capital"
	StopSpotWorkerCmd         = "post!/spot/stop"
	ArchiveSpotTradingDataCmd = "post!/spot/archive-data"

	GetWavetrendDataCmd = "get!/wavetrend-data"
)

func (t *teleBot) setupRoute() {
	// spot
	t.m[GetSpotAccountInfoCmd] = t.getSpotAccountInfo
	t.m[GetSpotHealthCmd] = t.healthCheckSpot

	t.m[CreateSpotWorkerCmd] = t.createNewSpotWorker
	t.m[AddCapitalSpotWorkerCmd] = t.addCapitalSpotWorker
	t.m[StopSpotWorkerCmd] = t.stopSpotWorker
	t.m[ArchiveSpotTradingDataCmd] = t.archiveSpotTradingData

	// wavetrend
	t.m[GetWavetrendDataCmd] = t.getWavetrendData

	t.bot.Handle(tb.OnText, func(c tb.Context) error {
		if c.Sender().ID != config.Instance().Common.AdminTeleID {
			c.Send("unauthorized 👾")
			return nil
		}

		args := strings.Fields(c.Text())
		cmd := args[0]
		handler, ok := t.m[cmd]
		if !ok {
			msg := "not found 🤖"
			c.Send(msg)
			return nil
		}
		handler(c)
		return nil
	})
}
