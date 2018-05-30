package rpsapp

import (
	"github.com/julienschmidt/httprouter"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo-games/rock-paper-scissors/server-go/rspbot"
)

func InitApp(httpRouter *httprouter.Router, botHost bots.BotHost) {
	rpsbot.InitBot(httpRouter, botHost, TheAppContext)
}
