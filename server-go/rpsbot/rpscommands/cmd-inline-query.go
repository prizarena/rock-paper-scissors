package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/bots-api-telegram"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"github.com/strongo/app"
	"github.com/prizarena/turn-based"
	"strings"
	"net/url"
)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		inlineQuery := whc.Input().(telegram.TgWebhookInlineQuery)
		inlineQueryID := inlineQuery.GetInlineQueryID()
		inlineQueryText := inlineQuery.TgUpdate().InlineQuery.Query

		switch {
		case strings.HasPrefix(inlineQueryText, "share?tournament="):
			// share tournament
		case strings.HasPrefix(inlineQueryText, "play"):
			return inlineQueryPlay(whc, inlineQueryID, inlineQueryText)
		}
		return
	},
)

func inlineQueryPlay(whc bots.WebhookContext, inlineQueryID, inlineQueryText string) (m bots.MessageFromBot, err error) {
	c := whc.Context()

	var values url.Values

	if qIndex := strings.Index(inlineQueryText, "?"); qIndex >= 0 {
		inlineQueryText = inlineQueryText[qIndex:]
		if values, err = url.ParseQuery(inlineQueryText); err != nil {
			return
		}
	}
	var tournamentID string
	if values != nil {
		tournamentID = values.Get("tournament")
	}

	translator := whc.BotAppContext().GetTranslator(c)
	newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
		t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)
		m, err := renderGameMessage(whc, t, turnbased.Board{BoardEntity: &turnbased.BoardEntity{TournamentID: tournamentID, Lang: lang, Round: 1}})
		if err != nil {
			panic(err)
		}
		return tgbotapi.InlineQueryResultArticle{
			ID:          "new-game_" + lang,
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
}
