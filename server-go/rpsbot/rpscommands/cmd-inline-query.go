package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpstrans"
	"github.com/strongo/app"
	)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		inlineQueryID := whc.Input().(telegram.TgWebhookInlineQuery).GetInlineQueryID()
		c := whc.Context()

		translator := whc.BotAppContext().GetTranslator(c)
		newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
			t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)
			callbackData := func(optionID string ) string {
				return "bet?on=" + t.Translate(optionID)
			}
			return tgbotapi.InlineQueryResultArticle{
				ID:          "new-game_"+lang,
				Type:        "article",
				Title:       t.Translate(rpstrans.NewGameInlineTitle),
				Description: t.Translate(rpstrans.NewGameInlineDescription),
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:                  t.Translate(rpstrans.NewGameText),
					ParseMode:             "HTML",
					DisableWebPagePreview: true,
				},
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					[]tgbotapi.InlineKeyboardButton{
						{Text: t.Translate(rpstrans.Option1text), CallbackData: callbackData(rpstrans.Option1code)},
					},
					[]tgbotapi.InlineKeyboardButton{
						{Text: t.Translate(rpstrans.Option2text), CallbackData: callbackData(rpstrans.Option2code)},
					},
					[]tgbotapi.InlineKeyboardButton{
						{Text: t.Translate(rpstrans.Option3text), CallbackData: callbackData(rpstrans.Option3code)},
					},
				),
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

