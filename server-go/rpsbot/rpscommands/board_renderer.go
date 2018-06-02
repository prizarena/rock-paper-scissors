package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/app"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsmodels"
	"bytes"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"fmt"
	"time"
	"strconv"
	"net/url"
	"github.com/strongo/bots-api-telegram"
)

func renderGameMessage(whc bots.WebhookContext, t strongo.SingleLocaleTranslator, lang, boardID string, round int, game rpsmodels.RpsGame, user rpsmodels.User) (m bots.MessageFromBot, err error) {
	m.IsEdit = true
	m.Format = bots.MessageFormatHTML
	var s bytes.Buffer
	s.WriteString(whc.Translate(rpstrans.GameCardTitle))
	s.WriteString("\n")
	if game.RpsGameEntity == nil {
		game.RpsGameEntity = &rpsmodels.RpsGameEntity{}
	}
	userMoves := game.UserMoves.Strings()
	movesCount := len(userMoves)
	switch {
	case movesCount == 0: // not started
		s.WriteString(whc.Translate(rpstrans.AskToMakeMove))
	case movesCount == 1 || userMoves[0] == "" || userMoves[1] == "": // one player made move, waiting for second player
		s.WriteString("\n\n")
		s.WriteString(t.Translate(rpstrans.FirstMoveDoneAwaitingSecond, user.ID))
	case movesCount == 2: // game completed
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Game completed, moves: %v v %v", userMoves[0], userMoves[1]))
		s.WriteString("\n\n")
		s.WriteString(whc.Translate(rpstrans.AskToMakeMove))
	default:
		panic(fmt.Sprintf("len(game.UserMoves) > 2: %v", movesCount))
	}
	s.WriteString(fmt.Sprintf("\n\nLast updated: %v", time.Now()))
	m.Text = s.String()
	callbackPrefix := fmt.Sprintf("bet?l=%v&r=%v&", lang, strconv.Itoa(round))
	if boardID != "" {
		callbackPrefix += "b=" + url.QueryEscape(boardID) + "&"
	}

	inlineQuery := ""
	inlineMove := func(text, id string) tgbotapi.InlineKeyboardButton {
		return tgbotapi.InlineKeyboardButton{Text: whc.Translate(text), CallbackData: callbackPrefix + "m=" + t.Translate(id)}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			inlineMove(rpstrans.Option1emoji, rpstrans.Option1code),
			inlineMove(rpstrans.Option2emoji, rpstrans.Option2code),
			inlineMove(rpstrans.Option3emoji, rpstrans.Option3code),
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:              whc.Translate(rpstrans.ChallengeFriendCommandText),
				SwitchInlineQuery: &inlineQuery,
			},
		},
	)

	m.Keyboard = keyboard
	return
}
