package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/turn-based"
	"net/url"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsdal"
	"context"
	"github.com/strongo/db"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
)

var leaveTournamentCommand = bots.NewCallbackCommand(turnbased.LeaveTournamentCommandCode, func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
	var board turnbased.Board
	board.ID = callbackUrl.Query().Get("board")
	if board.ID, err = turnbased.GetBoardID(whc.Input(), board.ID); err != nil {
		return
	}

	c := whc.Context()
	err = rpsdal.DB.RunInTransaction(c, func(c context.Context) error {
		err = turnbased.LeaveTournamentAction(whc, rpssecrets.PrizarenaGameID, rpssecrets.PrizarenaToken, rpsdal.DB, board)
		return nil
	}, db.SingleGroupTransaction)
	if err != nil {
		return
	}

	return renderRpsBoardMessage(whc, nil, board)
})