package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/prizarena-public/pabot"
	"net/http"
	"github.com/prizarena/prizarena-public/prizarena-interfaces"
	"github.com/prizarena/prizarena-public/prizarena-client-go"
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

	pabot.InitPrizarenaBot(router, func(httpClient *http.Client) prizarena_interfaces.ApiClient {
		if httpClient == nil {
			panic("httpClient == nil")
		}
		return prizarena.NewHttpApiClient(httpClient, "", rpssecrets.RpsPrizarenaGameID, rpssecrets.RpsPrizarenaToken)
	})
}
