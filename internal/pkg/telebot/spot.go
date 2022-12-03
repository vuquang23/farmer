package telebot

import (
	goctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/utils/context"
	"farmer/internal/pkg/utils/logger"
)

func (t *teleBot) getSpotAccountInfo(c tb.Context) {
	f := func(ctx goctx.Context) string {
		ret, err := t.spotTradeSvc.GetTradingPairsInfo(ctx)
		if err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		dtoRes := toGetSpotAccountInfoResponse(ret)
		return message(dtoRes)
	}

	msg := f(goctx.Background())
	c.Send(msg)
}

func toGetSpotAccountInfoResponse(en []*entities.SpotTradingPairInfo) *GetSpotAccountInfoResponse {
	var (
		p               = []*SpotPairInfo{}
		totalChangedUSD = 0.0
		totalBenefitUSD = 0.0
		N               = 10000.
	)

	for _, e := range en {
		info := SpotPairInfo{
			Symbol:          e.Symbol,
			Capital:         e.Capital,
			CurrentUSDValue: math.Round(e.CurrentUSDValue*N) / N,
			BenefitUSD:      math.Round(e.BenefitUSD*N) / N,
			ChangedUSD:      math.Round((e.CurrentUSDValue-e.Capital)*N) / N,
			BaseAmount:      math.Round(e.BaseAmount*N) / N,
			QuoteAmount:     math.Round(e.QuoteAmount*N) / N,
			UnitBuyAllowed:  e.UnitBuyAllowed,
			UnitNotional:    math.Round(e.UnitNotional*N) / N,
			TotalUnitBought: e.TotalUnitBought,
		}
		p = append(p, &info)

		totalBenefitUSD += info.BenefitUSD
		totalChangedUSD += info.ChangedUSD
	}

	return &GetSpotAccountInfoResponse{
		Pairs:           p,
		TotalBenefitUSD: math.Round(totalBenefitUSD*N) / N,
		TotalChangedUSD: math.Round(totalChangedUSD*N) / N,
	}
}

func (t *teleBot) healthCheckSpot(c tb.Context) {
	mapping := t.spotManager.CheckHealth()
	c.Send(message(mapping))
}

func (t *teleBot) createNewSpotWorker(c tb.Context) {
	f := func(ctx goctx.Context) string {
		args := strings.Fields(c.Text())
		if len(args) == 1 {
			return "missing required body"
		}

		var req CreateNewSpotWorkerReq
		if err := json.Unmarshal([]byte(args[1]), &req); err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		params := req.Normalize().ToCreateNewSpotWorkerParams()
		if err := t.spotManager.CreateNewWorker(
			context.Child(ctx, fmt.Sprintf("[create-new-worker] %s", params.Symbol)),
			params,
		); err != nil {
			return message(err)
		}

		return message("ok")
	}

	msg := f(goctx.Background())
	c.Send(msg)
}

func (t *teleBot) stopSpotWorker(c tb.Context) {
	f := func(ctx goctx.Context) string {
		args := strings.Fields(c.Text())
		if len(args) == 1 {
			return "missing required body"
		}

		var req StopWorkerReq
		if err := json.Unmarshal([]byte(args[1]), &req); err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		params := req.Normalize().ToStopBotParams()
		if err := t.spotManager.StopWorker(
			context.Child(ctx, fmt.Sprintf("[stop-bot] %s", params.Symbol)),
			params,
		); err != nil {
			return message(err)
		}

		return message("ok")
	}

	msg := f(goctx.Background())
	c.Send(msg)
}

func (t *teleBot) addCapitalSpotWorker(c tb.Context) {
	f := func(ctx goctx.Context) string {
		args := strings.Fields(c.Text())
		if len(args) == 1 {
			return "missing required body"
		}

		var req AddCapitalReq
		if err := json.Unmarshal([]byte(args[1]), &req); err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		params := req.Normalize().ToAddCapitalParams()
		if err := t.spotManager.AddCapital(
			context.Child(ctx, fmt.Sprintf("[add-capital-spot-worker] %s", params.Symbol)),
			params,
		); err != nil {
			return message(err)
		}

		return message("ok")
	}

	msg := f(goctx.Background())
	c.Send(msg)
}

func (t *teleBot) archiveSpotTradingData(c tb.Context) {
	f := func(ctx goctx.Context) string {

		args := strings.Fields(c.Text())
		if len(args) == 1 {
			return "missing required body"
		}

		var req ArchiveSpotTradingDataReq
		if err := json.Unmarshal([]byte(args[1]), &req); err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		params := req.Normalize().ToArchiveSpotTradingDataParams()
		ctx = context.Child(ctx, fmt.Sprintf("[archive-spot-trade-data] %s", params.Symbol))

		isActive := t.spotManager.IsActiveWorker(ctx, params.Symbol)
		if isActive {
			return message(errors.New("worker is active"))
		}

		err := t.spotTradeSvc.ArchiveTradingData(ctx, params)
		if err != nil {
			return message(err)
		}

		return message("ok")
	}

	msg := f(goctx.Background())
	c.Send(msg)
}
