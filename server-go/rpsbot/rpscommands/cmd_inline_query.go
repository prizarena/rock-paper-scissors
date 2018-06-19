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
			Text: tgInlineQuery.TgUpdate().InlineQuery.Query,
		}

		switch {
		case strings.HasPrefix(inlineQuery.Text, "tournament?id="):
			return inlineQueryTournament(whc, inlineQuery)
		case inlineQuery.Text == "play" || strings.HasPrefix(inlineQuery.Text, "play?tournament="):
			return inlineQueryPlay(whc, inlineQuery)
		default:
			return inlineQueryDefault(whc, inlineQuery)
		}
	},
)

func inlineQueryDefault(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	return
}

func inlineQueryTournament(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "inlineQuery: %+v", inlineQuery)
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery,  rpssecrets.RpsPrizarenaGameID, "id",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			log.Debugf(c,"inlineQueryTournament => pabot.ProcessInlineQueryTournament => reply")
			m, err = renderGameMessage(whc, whc, turnbased.Board{BoardEntity: &turnbased.BoardEntity{TournamentID: tournament.ShortTournamentID(), Lang: "en-US", Round: 1}})
			if err != nil {
				panic(err)
			}
			m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					tgbotapi.InlineQueryResultArticle{
						ID:          "newgame@t=" + tournament.ID,
						Type:        "article",
						Title:       "Start new game",
						Description: "At tournament: " + tournament.Name,
						InputMessageContent: tgbotapi.InputTextMessageContent{
							Text:                  m.Text,
							ParseMode:             "HTML",
							DisableWebPagePreview: m.DisableWebPagePreview,
						},
						ReplyMarkup: m.Keyboard.(*tgbotapi.InlineKeyboardMarkup),
					},
					//newGameOption("ru-RU"),
				},
			})
			m.Text = ""
			m.Keyboard = nil
			log.Debugf(c,"m: %+v", m)
			return
		})
}

func inlineQueryPlay(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery, rpssecrets.RpsPrizarenaGameID, "tournament",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			c := whc.Context()

			translator := whc.BotAppContext().GetTranslator(c)
			newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
				t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)
				m, err := renderGameMessage(whc, t, turnbased.Board{BoardEntity: &turnbased.BoardEntity{TournamentID: tournament.ID, Lang: lang, Round: 1}})
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
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					newGameOption("en-US"),
					//newGameOption("ru-RU"),
				},
			})
			return
		})
	return
}
