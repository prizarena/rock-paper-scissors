package rpscommands

import (
	"github.com/prizarena/turn-based"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsdal"
)

var leaveTournamentCommand = turnbased.NewLeaveTournamentCommand(
	rpssecrets.RpsPrizarenaGameID, rpssecrets.RpsPrizarenaToken, rpsdal.DB,
	func(whc bots.WebhookContext, board turnbased.Board) (m bots.MessageFromBot, err error) {
		return renderGameMessage(whc, whc, board)
	},
)
