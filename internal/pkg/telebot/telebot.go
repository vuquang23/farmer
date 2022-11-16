package telebot

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/pretty"
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/services"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/utils/logger"
)

type ITeleBot interface {
	Run(ctx context.Context)
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
func (tlb *teleBot) Run(ctx context.Context) {
	logger.Info(ctx, "[Run] start telebot")
	tlb.bot.Start()
}

// SendMsg sending a message to master
func (ltb *teleBot) SendMsg(ctx context.Context, msg interface{}) {
	for i := 0; i < 3; i++ {
		_, err := ltb.bot.Send(ltb.group, msg)
		if err != nil {
			logger.Info(ctx, "[SendMsg] %s", err)
			time.Sleep(3 * time.Second)
			continue
		}
		return
	}
}

// SendMsgWithFormat sending a message with format to master
func (ltb *teleBot) SendMsgWithFormat(ctx context.Context, template string, params ...interface{}) {
	ltb.SendMsg(ctx, fmt.Sprintf(template, params...))
}

func (t *teleBot) setupRoute() {
	t.m["get!/spot/account-info"] = t.getSpotAccountInfo
	t.m["get!/health"] = t.healthCheck

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
		ctx := context.Background()
		ret, err := t.spotTradeSvc.GetTradingPairsInfo(ctx)
		if err != nil {
			logger.Error(ctx, err)
			byteRes, _ := json.Marshal(err)
			return string(pretty.Pretty(byteRes))
		}

		dtoRes := toGetSpotAccountInfoResponse(ret)
		byteRes, _ := json.Marshal(dtoRes)
		return string(pretty.Pretty(byteRes))
	}

	msg := f()
	c.Send(msg)
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

func (t *teleBot) healthCheck(ctx tb.Context) {
	mapping := spotmanager.SpotManagerInstance().CheckHealth()
	bRes, _ := json.Marshal(mapping)
	response := string(pretty.Pretty(bRes))
	ctx.Send(response)
}
