package telebot

import (
	"encoding/json"
	"strings"

	tb "gopkg.in/telebot.v3"

	"farmer/internal/pkg/utils/logger"
)

func setupRoute(telebot *teleBot) {
	log := logger.WithDescription("Telegram bot controller")

	telebot.bot.Handle(tb.OnText, func(ctx tb.Context) error {
		args := strings.Fields(ctx.Text())
		switch args[0] {
		case "get!/spot/account-info":
			dto := &AccountInfoQuery{}
			if err := json.Unmarshal([]byte(args[1]), dto); err != nil {
				msg := "[get!/spotaccount-info] unmarshal error"
				log.Error(msg)
				ctx.Send(msg)
				return nil
			}
			getSpotAccountInfo(ctx, dto)

		case "post!/spot/bot": // create bot that will trade a SYMBOL
			req := &CreateSpotBotRequest{}
			if err := json.Unmarshal([]byte(args[1]), req); err != nil {
				msg := "[post!/spot/bot] unmarshal error"
				log.Error(msg)
				ctx.Send(msg)
				return nil
			}
			createSpotBot(ctx, req)

		case "delete!/spot/bot": // stop a bot that is trading a SYMBOL
			req := &StopBotBotRequest{}
			if err := json.Unmarshal([]byte(args[1]), req); err != nil {
				msg := "[delete!/spot/bot] unmarshal error"
				log.Error(msg)
				ctx.Send(msg)
				return nil
			}
			stopSpotBot(ctx, req)

		default:
			msg := "not match any route"
			log.Error(msg)
			ctx.Send(msg)
			return nil
		}
		return nil
	})
}
