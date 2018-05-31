package rpstgbots

import (
	"github.com/strongo/app"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpssecrets"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"context"
	"github.com/strongo/log"
)

var botsBy bots.SettingsBy

const BotProfile = "RockPaperScissors"

func Bots(c context.Context, env strongo.Environment, router bots.WebhooksRouter) bots.SettingsBy {
	if len(botsBy.ByCode) == 0 {
		routerByProfile := func(profile string) bots.WebhooksRouter {
			return router // We have single profile for now
		}

		switch env {
		case strongo.EnvProduction:
			botsBy = bots.NewBotSettingsBy(routerByProfile,
				telegram.NewTelegramBot(strongo.EnvProduction, BotProfile,
					rpssecrets.TelegramProdBot, rpssecrets.TelegramProdToken,
					"", "", rpssecrets.GaTrackingID, strongo.LocaleEnUS),
			)
		//case strongo.EnvLocal:
		//	botsBy = bots.NewBotSettingsBy(routerByProfile,
		//		telegram.NewTelegramBot(strongo.EnvLocal, BotProfile,
		//			rpssecrets.TelegramLocalBot, rpssecrets.TelegramLocalToken,
		//			"", "", "", strongo.LocaleEnUS),
		//	)
		default:
			log.Errorf(c, "Unknown environment: %v=%v", env, strongo.EnvironmentNames[env])
		}
	}
	return botsBy
}

func BotsBy(_ context.Context) bots.SettingsBy {
	return botsBy
}