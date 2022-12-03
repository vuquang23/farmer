package telebot

import (
	goctx "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	"github.com/tidwall/pretty"
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/services"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/utils/context"
	"farmer/internal/pkg/utils/logger"
	wavetrendprovider "farmer/internal/pkg/wavetrend"
)

type ITeleBot interface {
	Run(ctx goctx.Context)
	SendMsg(ctx context.Context, msg interface{})
	SendMsgWithFormat(ctx context.Context, template string, params ...interface{})
}

type handlerFunc func(c tb.Context)

// TeleBot ...
type teleBot struct {
	bot               *tb.Bot
	group             *tb.Chat
	spotTradeSvc      services.ISpotTradeService
	wavetrendProvider wavetrendprovider.IWavetrendProvider
	spotManager       spotmanager.ISpotManager
	m                 map[string]handlerFunc
}

var tlbot *teleBot

func TeleBotInstance() ITeleBot {
	return tlbot
}

func InitTeleBot(spotTradeSvc services.ISpotTradeService, wavetrendProvider wavetrendprovider.IWavetrendProvider, spotManager spotmanager.ISpotManager) error {
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
		spotTradeSvc:      spotTradeSvc,
		wavetrendProvider: wavetrendProvider,
		spotManager:       spotManager,
		m:                 make(map[string]handlerFunc),
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

func message(data interface{}) string {
	dataStr, ok := data.(string)
	if ok {
		return dataStr
	}
	dataBytes, _ := json.Marshal(data)
	return string(pretty.Pretty(dataBytes))
}
