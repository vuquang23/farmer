package telebot

import (
	goctx "context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	"github.com/tidwall/pretty"
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/services"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/utils/context"
	"farmer/internal/pkg/utils/logger"
)

type ITeleBot interface {
	Run(ctx goctx.Context)
	SendMsg(ctx context.Context, msg interface{})
	SendMsgWithFormat(ctx context.Context, template string, params ...interface{})
}

type handlerFunc func(c tb.Context)

// TeleBot ...
type teleBot struct {
	bot          *tb.Bot
	group        *tb.Chat
	spotTradeSvc services.ISpotTradeService
	m            map[string]handlerFunc
}

var tlbot *teleBot

func TeleBotInstance() ITeleBot {
	return tlbot
}

func InitTeleBot(spotTradeSvc services.ISpotTradeService) error {
	if tlbot != nil {
		return nil
	}

	token := config.Instance().Telebot.Token
	groupID := int64(config.Instance().Telebot.GroupID)
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return errors.Wrap(err, "failed to create bot")
	}

	tlbot = &teleBot{
		bot: bot,
		group: &tb.Chat{
			ID:   groupID,
			Type: tb.ChatGroup,
		},
		spotTradeSvc: spotTradeSvc,
		m:            make(map[string]handlerFunc),
	}
	tlbot.setupRoute()

	return nil
}

// Run runs bot
func (t *teleBot) Run(ctx goctx.Context) {
	logger.Info(ctx, "[Run] start telebot")
	t.bot.Start()
}

// SendMsg sending a message to master
func (ltb *teleBot) SendMsg(ctx context.Context, msg interface{}) {
	var (
		tried   = 0
		backoff = retry.NewFibonacci(1 * time.Second)
	)

	_ = retry.Do(ctx, retry.WithMaxRetries(2, backoff), func(ctx goctx.Context) error {
		defer func() {
			tried++
		}()
		if tried > 0 {
			logger.Infof(ctx, "[SendMsg] retry %d ...", tried)
		}

		_, err := ltb.bot.Send(ltb.group, msg)
		if err != nil {
			logger.Error(ctx, "[SendMsg] %s", err)
			return err
		}

		return nil
	})
}

// SendMsgWithFormat sending a message with format to master
func (t *teleBot) SendMsgWithFormat(ctx context.Context, template string, params ...interface{}) {
	t.SendMsg(ctx, fmt.Sprintf(template, params...))
}

func (t *teleBot) setupRoute() {
	t.m["get!/spot/account-info"] = t.getSpotAccountInfo
	t.m["get!/health"] = t.healthCheck

	t.m["post!/spot"] = t.createNewSpotWorker
	t.m["post!/spot/add-capital"] = t.addCapitalSpotWorker
	t.m["post!/spot/stop"] = t.stopSpotBot

	t.bot.Handle(tb.OnText, func(c tb.Context) error {
		args := strings.Fields(c.Text())
		cmd := args[0]
		handler, ok := t.m[cmd]
		if !ok {
			msg := "not found"
			c.Send(msg)
			return nil
		}
		handler(c)
		return nil
	})
}

func (t *teleBot) getSpotAccountInfo(c tb.Context) {
	f := func() string {
		ctx := goctx.Background()
		ret, err := t.spotTradeSvc.GetTradingPairsInfo(ctx)
		if err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		dtoRes := toGetSpotAccountInfoResponse(ret)
		return message(dtoRes)
	}

	msg := f()
	c.Send(msg)
}

func toGetSpotAccountInfoResponse(en []*entities.SpotTradingPairInfo) *GetSpotAccountInfoResponse {
	p := []*SpotPairInfo{}
	totalBenefitUSD := 0.0
	N := 10000.

	for _, e := range en {
		p = append(p, &SpotPairInfo{
			Symbol:          e.Symbol,
			Capital:         e.Capital,
			CurrentUSDValue: math.Round(e.CurrentUSDValue*N) / N,
			BenefitUSD:      math.Round(e.BenefitUSD*N) / N,
			BaseAmount:      math.Round(e.BaseAmount*N) / N,
			QuoteAmount:     math.Round(e.QuoteAmount*N) / N,
			UnitBuyAllowed:  e.UnitBuyAllowed,
			UnitNotional:    math.Round(e.UnitNotional*N) / N,
			TotalUnitBought: e.TotalUnitBought,
		})

		totalBenefitUSD += e.BenefitUSD
	}

	return &GetSpotAccountInfoResponse{
		Pairs:           p,
		TotalBenefitUSD: math.Round(totalBenefitUSD*N) / N,
	}
}

func (t *teleBot) healthCheck(c tb.Context) {
	mapping := spotmanager.SpotManagerInstance().CheckHealth()
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
		if err := spotmanager.SpotManagerInstance().CreateNewWorker(
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

func (t *teleBot) stopSpotBot(c tb.Context) {
	f := func(ctx goctx.Context) string {
		args := strings.Fields(c.Text())
		if len(args) == 1 {
			return "missing required body"
		}

		var req StopBotReq
		if err := json.Unmarshal([]byte(args[1]), &req); err != nil {
			logger.Error(ctx, err)
			return message(err)
		}

		params := req.Normalize().ToStopBotParams()
		if err := spotmanager.SpotManagerInstance().StopBot(
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
		if err := spotmanager.SpotManagerInstance().AddCapital(
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

func message(data interface{}) string {
	dataStr, ok := data.(string)
	if ok {
		return dataStr
	}
	dataBytes, _ := json.Marshal(data)
	return string(pretty.Pretty(dataBytes))
}
