package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/prizarena-public/pabot"
				"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
)

func RegisterRpsCommands(router bots.WebhooksRouter) {
	router.RegisterCommands([]bots.Command{
		startCommand,
		inlineQueryCommand,
		chosenInlineResultCommand,
		betCallbackCommand,
		leaveTournamentCommand,
	})

	pabot.InitPrizarenaInGameBot(rpssecrets.PrizarenaGameID, rpssecrets.PrizarenaToken, router)
}
