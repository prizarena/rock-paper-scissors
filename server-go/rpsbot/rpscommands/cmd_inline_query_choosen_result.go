package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"github.com/prizarena/turn-based"
	"strings"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsdal"
	"context"
	"github.com/strongo/db"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
	"time"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsmodels"
)

var chosenInlineResultCommand = bots.Command{
	Code:       "inline-chosen-result",
	InputTypes: []bots.WebhookInputType{bots.WebhookInputChosenInlineResult},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		chosenResult := whc.Input().(bots.WebhookChosenInlineResult)
		query := chosenResult.GetQuery()
		log.Infof(c, "chosenInlineResultCommand.Action(): query: %v", query)

		if strings.HasPrefix(query, "play") {
			var tournament pamodels.Tournament
			if tournament.ID, err = pabot.GetQueryValueFromUrl(query, "tournament"); err != nil {
				return
			}
			if tournament.ID != "" {
				tournament.ID = rpssecrets.PrizarenaGameID + pamodels.TournamentIDSeparator + tournament.ID
				prizarenaCachedApi := pabot.NewCachedApi(c, rpssecrets.PrizarenaGameID, rpssecrets.PrizarenaToken)
				if tournament, err = prizarenaCachedApi.GetTournamentByID(c, tournament.ID); err != nil {
					return
				}
			}

			var board turnbased.Board
			board.ID = chosenResult.GetInlineMessageID()
			if err = rpsdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
				if err = rpsdal.DB.Get(c, &board); err == nil {
					// Board entity already created
					if board.TournamentID != tournament.ID {
						log.Warningf(c, "board.ID=%v => board.TournamentID != tournament.ID => %v != %v ", board.ID, board.TournamentID, tournament.ID)
					}
					return
				} else if !db.IsNotFound(err) {
					return
				}
				var botAppUser bots.BotAppUser
				if botAppUser, err = whc.GetAppUser(); err != nil {
					return
				}
				appUserEntity := botAppUser.(*rpsmodels.UserEntity)

				board.BoardEntity = &turnbased.BoardEntity{
					BoardEntityBase: turnbased.BoardEntityBase{
						Round:        1,
						TournamentID: tournament.ShortTournamentID(),
						UserIDs:      []string{whc.AppUserStrID()},
						UserNames:    []string{appUserEntity.FullName()},
						Created:      time.Now(),
					},
				}
				m, err = renderRpsBoardMessage(whc, &tournament, board)
				return
			}, db.CrossGroupTransaction); err != nil {
				return
			}
		}
		return
	},
}
