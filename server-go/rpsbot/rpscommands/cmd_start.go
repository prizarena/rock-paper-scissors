package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
)

var startCommand = bots.Command{
	Code:     "start",
	Commands: []string{"/start"},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		if m, err = pabot.OnStartIfTournamentLink(whc, rpssecrets.PrizarenaGameID); m.Text != "" || err != nil {
			return
		}
		m.Text = whc.Translate(rpstrans.NewGameText, whc.Translate(rpstrans.RulesShort))
		m.Format = bots.MessageFormatHTML
		m.DisableWebPagePreview = true
		inlineQuery := ""
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: whc.Translate(rpstrans.ChallengeFriendCommandText), SwitchInlineQuery: &inlineQuery},
			},
			pabot.NewTournamentTelegramInlineButton(whc, "rockpaperscissors"),
			[]tgbotapi.InlineKeyboardButton{
				{Text: "📜 Rules & How to play", URL: "https://prizarena.com/rockpaperscissors/"},
			},
		)
		return
	},
}
