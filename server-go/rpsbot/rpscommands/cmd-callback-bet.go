package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsdal"
	"context"
	"github.com/strongo/db"
	"github.com/prizarena/turn-based"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsmodels"
	"github.com/strongo/log"
	"github.com/pkg/errors"
	"time"
	"strconv"
)

var betCallbackCommand = bots.NewCallbackCommand(
	"bet",
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		q := callbackUrl.Query()
		lang := q.Get("l")
		whc.SetLocale(lang)
		boardID := q.Get("b")
		move := q.Get("m")
		var round int
		round, err = strconv.Atoi(q.Get("r"))
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
		var board turnbased.Board

		if err = rpsdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			board, err = turnbased.MakeMove(tc, time.Now(), rpsdal.DB, round, lang, boardID, userID, move)
			if err != nil {
				if err == turnbased.ErrOldRound || err == turnbased.ErrUnknownRound {
					log.Warningf(c, "%v: %v", err, round)
					err = nil
					return
				}
				return
			}
			rpsGame = rpsmodels.RpsGame{
				RpsGameEntity: &rpsmodels.RpsGameEntity{
					BoardID:     boardID,
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
				turnbased.NextRound(board)
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
			IntegerID:  db.NewIntID(whc.AppUserIntID()),
			UserEntity: appUser.(*rpsmodels.UserEntity),
		}
		if m, err = renderGameMessage(whc, whc, lang, boardID, board.Round, rpsGame, user); err != nil {
			return
		}
		return
	},
)
