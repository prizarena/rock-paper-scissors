package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-api-telegram"
)

var startCommand = bots.Command{
	Code: "start",
	Commands: []string{"/start"},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		m.Text = "Rock - Paper - Scissors"
		inlineQuery := ""
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: "Challenge Telegram friend", SwitchInlineQuery: &inlineQuery},
			},
			[]tgbotapi.InlineKeyboardButton{
				{Text: "New tournament", CallbackData: "new-tournament"},
			},
			[]tgbotapi.InlineKeyboardButton{
				{Text: "Rules & How to play", CallbackData: "rules"},
			},
		)
		return
	},
}
