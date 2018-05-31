package rpsapp

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsbot"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func InitApp(botHost bots.BotHost) {
	httpRouter := httprouter.New()
	http.Handle("/", httpRouter)

	httpRouter.GET("/test1", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte("test response 1"))
	})
	httpRouter.GET("/test2", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte("test response 2"))
	})

	rpsbot.InitBot(httpRouter, botHost, rpsAppContext{})
}
