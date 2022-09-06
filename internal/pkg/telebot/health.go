package telebot

import (
	"encoding/json"

	"github.com/tidwall/pretty"
	tb "gopkg.in/telebot.v3"

	spotmanager "farmer/internal/pkg/spot_manager"
)

func (tlb *teleBot) checkHealth(ctx tb.Context) {
	mapping := spotmanager.SpotManagerInstance().CheckHealth()
	bRes, _ := json.Marshal(mapping)
	response := string(pretty.Pretty(bRes))
	ctx.Send(response)
}
