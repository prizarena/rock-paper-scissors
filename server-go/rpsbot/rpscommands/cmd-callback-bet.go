package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsdal"
	"context"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsmodels"
	"github.com/strongo/db"
)

var betCallbackCommand = bots.NewCallbackCommand(
	"bet",
	func(whc bots.WebhookContext, callbackUrl *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		bet := callbackUrl.Query().Get("on")

		if err = rpsdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
			canvas := rpsmodels.GameCanvas{StringID: db.NewStrID("")}
			err = rpsdal.DB.Get(c, &canvas)

			switch db.IsNotFound(err) {
			case true: // New canvas
				canvas.GameCanvasEntity = &rpsmodels.GameCanvasEntity{
					UserIDs: []string{whc.AppUserStrID(), ""},
				}
				if err = rpsdal.DB.Update(c, &canvas); err != nil {
					return
				}
			case false:
				if err != nil {
					return
				}

			}
			return
		}, db.SingleGroupTransaction); err != nil {
			return
		}

		m.IsEdit = true
		m.Text = "<b>Selected:</b> " + bet
		return
	},
)
