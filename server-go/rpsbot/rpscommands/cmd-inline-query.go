package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/bots-api-telegram"
)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		inlineQueryID := whc.Input().(telegram.TgWebhookInlineQuery).GetInlineQueryID()

		inlineResult := func(id, title, desc string) tgbotapi.InlineQueryResultArticle {
			return tgbotapi.InlineQueryResultArticle{
				ID:          id,
				Type:        "article",
				Title:       title,
				Description: desc,
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:                  `<b>Rock-Paper-Scissors</b>

Please make your choice.

<code>Rules</code>:
<pre>
  💎 Rock wins over scissors ✂.

  📄 Paper wins over rock 💎.

  ✂ Scissors win over paper 📄.
</pre>
<b>Sponsored:</b> <a href="https://t.me/playRockPaperScissorsBot?start=place-ad">Your ad could be here.</a>`,
					ParseMode:             "HTML",
					DisableWebPagePreview: true,
				},
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					[]tgbotapi.InlineKeyboardButton{
						{Text: "💎 Rock", CallbackData: "rock"},
					},
					[]tgbotapi.InlineKeyboardButton{
						{Text: "📄 Paper", CallbackData: "paper"},
					},
					[]tgbotapi.InlineKeyboardButton{
						{Text: "✂ Scissors", CallbackData: "scissors"},
					},
				),
			}
		}
		m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
			InlineQueryID: inlineQueryID,
			Results: []interface{}{
				inlineResult("new-game", "💎📄✂ New game ", "Starts new Rock-Paper-Scissors game"),
				//inlineResult("paper", "📄 Paper", "Wins over rock 💎"),
				//inlineResult("scissors", "✂ Scissors", "Wins over paper 📄"),
			},
		})
		return
	},
)

