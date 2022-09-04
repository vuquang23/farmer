package telebot

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/utils/logger"
)

type ITeleBot interface {
	Run()
	SendMsg(msg interface{})
	SendMsgWithFormat(template string, params ...interface{})
}

// TeleBot ...
type teleBot struct {
	bot   *tb.Bot
	group *tb.Chat

	spotTradeSvc services.ISpotTradeService
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
	}

	setupRoute(tlbot)

	return nil
}

// Run run bot
func (tlb *teleBot) Run() {
	tlb.bot.Start()
}

// SendMsg sending a message to master
func (ltb *teleBot) SendMsg(msg interface{}) {
	log := logger.WithDescription("Send message telegram")

	for i := 0; i < 3; i++ {
		_, err := ltb.bot.Send(ltb.group, msg)
		if err != nil {
			log.Sugar().Info("failed to send message tele - ", "attempt", i, " - err", err)
			time.Sleep(3 * time.Second)
			continue
		}
		return
	}
}

// SendMsgWithFormat sending a message with format to master
func (ltb *teleBot) SendMsgWithFormat(template string, params ...interface{}) {
	ltb.SendMsg(fmt.Sprintf(template, params...))
}
