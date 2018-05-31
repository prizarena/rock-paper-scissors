package rpsapp

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsbot"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsdal/rpsgaedal"
)

func InitApp(botHost bots.BotHost) {
	rpsgaedal.RegisterDal()

	httpRouter := httprouter.New()
	http.Handle("/", httpRouter)

	rpsbot.InitBot(httpRouter, botHost, rpsAppContext{})
}
