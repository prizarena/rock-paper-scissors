package rpsgaeroot

import (
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsapp"
	"github.com/strongo/log"
	"github.com/strongo/bots-framework/hosts/appengine"
	"github.com/julienschmidt/httprouter"
)

func init() {
	log.AddLogger(gaehost.GaeLogger)

	httpRouter := httprouter.New()

	rpsapp.InitApp(httpRouter, gaehost.GaeBotHost{})
}
