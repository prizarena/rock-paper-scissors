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
			"github.com/prizarena/prizarena-public/prizarena-interfaces"
	"github.com/prizarena/prizarena-public/pabot"
)

func getBoardID(whc bots.WebhookContext, boardID string) (string, error) {
	if boardID == "" {
		boardID = whc.Input().(bots.WebhookCallbackQuery).GetInlineMessageID()
		if boardID == "" {
			return "", errors.New("expecting to get inlineMessageID")
		}
	}
	return boardID, nil
}

var betCallbackCommand = bots.NewCallbackCommand(
	"bet",
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		q := callbackUrl.Query()
		lang := q.Get("l")
		if lang == "" {
			lang = "en-US"
		} else {
			whc.SetLocale(lang)
		}
		boardID := q.Get("b")
		move := q.Get("m")
		var round int
		round, err = strconv.Atoi(q.Get("r"))
		if err != nil {
			err = errors.WithMessage(err, "parameter 'round' is required and have to be an integer")
			return
		}

		if boardID, err = getBoardID(whc, boardID); err != nil {
			return
		}

		userID := whc.AppUserStrID()

		var board turnbased.Board

		var appUser bots.BotAppUser
		appUser, err = whc.GetAppUser()
		userEntity := appUser.(*rpsmodels.UserEntity)

		if err = rpsdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
			board, err = turnbased.MakeMove(tc, time.Now(), rpsdal.DB, round, lang, boardID, userID, userEntity.FullName(), move)
			if err != nil {
				if err == turnbased.ErrOldRound || err == turnbased.ErrUnknownRound {
					log.Warningf(c, "%v: %v", err, round)
					err = nil
					return
				}
				return
			}
			userMoves := board.UserMoves.Strings()
			if len(userMoves) == 2 && userMoves[0] != "" && userMoves[1] != "" {
				turnbased.NextRound(board)
			}
			if err = rpsdal.DB.Update(tc, &board); err != nil {
				return
			}
			return
		}, db.SingleGroupTransaction); err != nil {
			return
		}


		httpClient := whc.BotContext().BotHost.GetHTTPClient(c)
		prizarenaApiClient := pabot.NewPrizarenaApiClient(httpClient)

		prizarenaApiClient.PlayCompleted(c, prizarena_interfaces.PlayCompletedPayload{
			PlayID: "",
			TournamentID: "", // empty for monthly tournament
			Impacts: []prizarena_interfaces.Impact{
				{UserID: userID, Points: 1},
			},
		})
		if m, err = renderRpsBoardMessage(whc, nil, board); err != nil {
			return
		}
		return
	},
)
