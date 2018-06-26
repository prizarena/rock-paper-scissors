package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"strings"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/app"
	"github.com/prizarena/turn-based"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
	"github.com/strongo/log"
	)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		tgInlineQuery := whc.Input().(telegram.TgWebhookInlineQuery)
		inlineQuery := pabot.InlineQueryContext{
			ID:   tgInlineQuery.GetInlineQueryID(),
			Text: strings.TrimSpace(tgInlineQuery.TgUpdate().InlineQuery.Query),
		}

		switch {
		case strings.HasPrefix(inlineQuery.Text, "tournament?id="):
			return inlineQueryTournament(whc, inlineQuery)
		case inlineQuery.Text == "" || inlineQuery.Text == "play" || strings.HasPrefix(inlineQuery.Text, "play?tournament="):
			return inlineQueryPlay(whc, inlineQuery)
		}
		return
	},
)

// func inlineQueryDefault(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
// 	return
// }

func inlineQueryTournament(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "inlineQuery: %+v", inlineQuery)
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery, rpssecrets.PrizarenaGameID, rpssecrets.PrizarenaToken, "id",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			if tournament.TournamentEntity == nil {
				return
			}
			log.Debugf(c, "inlineQueryTournament => pabot.ProcessInlineQueryTournament => reply")
			newBoard := turnbased.Board{BoardEntity: &turnbased.BoardEntity{BoardEntityBase: turnbased.BoardEntityBase{Lang: "en-US", Round: 1}}}
			if m, err = renderRpsBoardMessage(whc, &tournament, newBoard); err != nil {
				panic(err)
			}
			var description string

			if tournament.IsSponsored {
				description = "Sponsored"
			} else {
				description = "Not sponsored"
			}
			m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					tgbotapi.InlineQueryResultArticle{
						ID:          "tournament=" + tournament.ID,
						Type:        "article",
						Title:       "💎📄✂ Tournament: " + tournament.Name,
						Description: description,
						InputMessageContent: tgbotapi.InputTextMessageContent{
							Text:                  m.Text,
							ParseMode:             "HTML",
							DisableWebPagePreview: m.DisableWebPagePreview,
						},
						ReplyMarkup: m.Keyboard.(*tgbotapi.InlineKeyboardMarkup),
					},
					// newGameOption("ru-RU"),
				},
			})
			m.Text = ""
			m.Keyboard = nil
			log.Debugf(c, "m: %+v", m)
			return
		})
}

func inlineQueryPlay(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery, rpssecrets.PrizarenaGameID, rpssecrets.PrizarenaToken, "tournament",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			c := whc.Context()

			translator := whc.BotAppContext().GetTranslator(c)

			newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
				t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)
				newBoard := turnbased.Board{BoardEntity: &turnbased.BoardEntity{BoardEntityBase: turnbased.BoardEntityBase{Lang: "en-US", Round: 1}}}

				// Renders game board to a Telegram message to return as inline result
				if m, err = renderRpsBoardMessage(t, &tournament, newBoard); err != nil {
					panic(err)
				}

				articleID := "new_game?l=" + lang
				if tournament.ID != "" {
					articleID += "&t=" + tournament.ShortTournamentID()
				}
				return tgbotapi.InlineQueryResultArticle{
					ID:          articleID,
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
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					newGameOption("en-US"),
					// newGameOption("ru-RU"),
				},
			})
			return
		})
	return
}
