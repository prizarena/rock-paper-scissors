package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsdal"
	"context"
	"github.com/strongo/db"
	"github.com/strongo-games/arena/canvas"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsmodels"
	"github.com/strongo/log"
	"bytes"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpstrans"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/app"
	"fmt"
	"github.com/pkg/errors"
	"time"
	"strconv"
)

var betCallbackCommand = bots.NewCallbackCommand(
	"bet",
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		q := callbackUrl.Query()
		boardID := q.Get("board")
		move :=  q.Get("move")
		var round int
		round, err = strconv.Atoi(q.Get("round"))
		if err != nil {
			err = errors.WithMessage(err, "parameter 'round' is required and have to be an integer")
			return
		}

		if boardID == "" {
			boardID = whc.Input().(bots.WebhookCallbackQuery).GetInlineMessageID()
			if boardID == "" {
				err = errors.New("expecting to get inlineMessageID")
			}
		}
		userID := whc.AppUserStrID()

		var rpsGame rpsmodels.RpsGame
		var board canvas.Board

		if err = rpsdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			board, err = canvas.MakeMove(tc, time.Now(), rpsdal.DB, round, boardID, userID, move)
			if err != nil {
				if err == canvas.ErrOldRound || err == canvas.ErrUnknownRound {
					log.Warningf(c, "%v: %v", err, round)
					err = nil
					return
				}
				return
			}
			rpsGame = rpsmodels.RpsGame{
				RpsGameEntity: &rpsmodels.RpsGameEntity{
					BoardID: boardID,
					BoardEntity: *board.BoardEntity,
				},
			}
			userMoves := board.UserMoves.Strings()
			if len(userMoves) == 2 && userMoves[0] != "" && userMoves[1] != "" {
				go func() {
					if err := rpsdal.DB.InsertWithRandomIntID(c, &rpsGame); err != nil {
						log.Errorf(c, "Failed to create RpsGame entity: %v", err)
					}
				}()
				canvas.NextRound(board)
			}
			if err = rpsdal.DB.Update(tc, &board); err != nil {
				return
			}
			return
		}, db.SingleGroupTransaction); err != nil {
			return
		}

		var appUser bots.BotAppUser
		appUser, err = whc.GetAppUser()
		if err != nil {
			return
		}
		user := rpsmodels.User{
			IntegerID: db.NewIntID(whc.AppUserIntID()),
			UserEntity: appUser.(*rpsmodels.UserEntity),
		}
		if m, err = renderGameMessage(whc, whc, boardID, board.Round, rpsGame, user); err != nil {
			return
		}
		return
	},
)

func renderGameMessage(whc bots.WebhookContext, t strongo.SingleLocaleTranslator, boardID string, round int, game rpsmodels.RpsGame, user rpsmodels.User) (m bots.MessageFromBot, err error) {
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
	callbackPrefix := "bet?round=" + strconv.Itoa(round) + "&"
	if boardID != "" {
		callbackPrefix += "board=" + url.QueryEscape(boardID) + "&"
	}

	callbackData := func(optionID string ) string {
		return callbackPrefix + "move=" + t.Translate(optionID)
	}
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: t.Translate(rpstrans.Option1text), CallbackData: callbackData(rpstrans.Option1code)},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: t.Translate(rpstrans.Option2text), CallbackData: callbackData(rpstrans.Option2code)},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: t.Translate(rpstrans.Option3text), CallbackData: callbackData(rpstrans.Option3code)},
		},
	)
	return
}