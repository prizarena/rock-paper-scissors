package rpsrouting

import (
	"github.com/strongo-games/bidding-tictactoe/server-go/btttbot/commands"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsbot/rpscommands"
)

var WebhooksRouter = bots.NewWebhookRouter(
	map[bots.WebhookInputType][]bots.Command{
		//bots.WebhookInputText: {
		//	commands.StartCommand,
		//	commands.BidCommand,
		//},
		bots.WebhookInputChosenInlineResult: {
			commands.ChosenInlineResultCommand,
		},
		bots.WebhookInputCallbackQuery: {
			commands.NewGameCommand,
			commands.HitCommand,
			commands.CreateBoardCommand,
			commands.StartLanguageSelectedCommand,
		},
	},
	func() string { return "Oops..." },
)

func init() {
	rpscommands.RegisterRpsCommands(WebhooksRouter)
}