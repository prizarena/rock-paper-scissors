package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
)

var chosenInlineResultCommand = bots.Command{
	Code:       "inline-choosen-result",
	InputTypes: []bots.WebhookInputType{bots.WebhookInputChosenInlineResult},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		chosenResult := whc.Input().(bots.WebhookChosenInlineResult)
		query := chosenResult.GetQuery()
		log.Infof(c, "chosenInlineResultCommand.Action(): query: %v", query)
		//if strings.HasPrefix(query, "#") {
		//	m, err = bidConfirmed(whc, query)
		//}
		return
	},
}
