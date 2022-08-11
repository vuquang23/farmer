package telebot

import (
	"farmer/internal/pkg/utils/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	tb "gopkg.in/telebot.v3"
)

// TeleBot ...
type TeleBot struct {
	ctx   *gin.Context
	bot   *tb.Bot
	group *tb.Chat
}

var telebot *TeleBot

func TeleBotInstance() *TeleBot {
	return telebot
}

func InitTeleBot(ctx *gin.Context, token string, groupID int64) error {
	if telebot != nil {
		return nil
	}
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return errors.Wrap(err, "failed to create bot")
	}
	telebot = &TeleBot{
		ctx: ctx,
		bot: bot,
		group: &tb.Chat{
			ID:   groupID,
			Type: tb.ChatGroup,
		},
	}

	setupRoute(telebot)

	return nil
}

// Run run bot
func (ltb *TeleBot) Run() {
	logger.FromGinCtx(ltb.ctx).Info("start telebot")
	ltb.bot.Start()
}

// SendMsg sending a message to master
func (ltb *TeleBot) SendMsg(msg interface{}) {
	l := logger.FromGinCtx(ltb.ctx)

	for i := 0; i < 3; i++ {
		_, err := ltb.bot.Send(ltb.group, msg)
		if err != nil {
			l.Sugar().Info("failed to send message tele - ", "attempt", i, " - err", err)
			time.Sleep(3 * time.Second)
			continue
		}
		return
	}
}

// SendMsgWithFormat sending a message with format to master
func (ltb *TeleBot) SendMsgWithFormat(template string, params ...interface{}) {
	ltb.SendMsg(fmt.Sprintf(template, params...))
}
