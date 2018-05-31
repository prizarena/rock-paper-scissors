package rpscommands

import "github.com/strongo/bots-framework/core"

func RegisterRpsCommands(router bots.WebhooksRouter) {
	router.RegisterCommands([]bots.Command{
		startCommand,
		inlineQueryCommand,
		chosenInlineResultCommand,
		betCallbackCommand,
	})
}
