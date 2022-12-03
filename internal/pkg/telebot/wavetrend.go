package telebot

import (
	goctx "context"
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/utils/context"
)

func (t *teleBot) getWavetrendData(c tb.Context) {
	f := func(ctx goctx.Context) string {
		args := strings.Fields(c.Text())
		if len(args) < 3 {
			return "missing svcName"
		}
		svcName := strings.ToUpper(args[1]) + ":" + strings.ToLower(args[2])

		var ret GetWavetrendDataResponse

		data, isOutdated := t.wavetrendProvider.GetPastWaveTrendData(
			context.Child(ctx, fmt.Sprintf("[get-wavetrend-data] %s", svcName)),
			svcName,
		)
		ret.IsOutdated = isOutdated
		ret.PastTci = data.PastTci
		ret.DifWavetrend = data.DifWavetrend

		currentTci, _ := t.wavetrendProvider.GetCurrentTci(ctx, svcName)
		ret.CurrentTci = currentTci

		currentDifWavetrend, _ := t.wavetrendProvider.GetCurrentDifWavetrend(ctx, svcName)
		ret.CurrentDifWavetrend = currentDifWavetrend

		closePrice, _ := t.wavetrendProvider.GetClosePrice(ctx, svcName)
		ret.ClosePrice = closePrice

		return message(ret)
	}

	msg := f(goctx.Background())
	c.Send(msg)
}
