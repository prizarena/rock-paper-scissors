package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/bots-api-telegram"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"github.com/strongo/app"
	"github.com/prizarena/turn-based"
)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		inlineQueryID := whc.Input().(telegram.TgWebhookInlineQuery).GetInlineQueryID()
		c := whc.Context()

		translator := whc.BotAppContext().GetTranslator(c)
		newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
			t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)
			m, err := renderGameMessage(whc, t, turnbased.Board{BoardEntity: &turnbased.BoardEntity{Lang: lang}})
			if  err != nil {
				panic(err)
			}
			return tgbotapi.InlineQueryResultArticle{
				ID:          "new-game_"+lang,
				Type:        "article",
				Title:       t.Translate(rpstrans.NewGameInlineTitle),
				Description: t.Translate(rpstrans.NewGameInlineDescription),
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:                  m.Text,
					ParseMode:             "HTML",
					DisableWebPagePreview: m.DisableWebPagePreview,
				},
				ReplyMarkup: m.Keyboard.(*tgbotapi.InlineKeyboardMarkup),
			}
		}

		m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
			InlineQueryID: inlineQueryID,
			Results: []interface{}{
				newGameOption("en-US"),
				newGameOption("ru-RU"),
			},
		})
		return
	},
)

