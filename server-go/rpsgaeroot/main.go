package rpsgaeroot

import (
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsapp"
	"github.com/strongo/log"
	"github.com/strongo/bots-framework/hosts/appengine"
)

func init() {
	log.AddLogger(gaehost.GaeLogger)
	rpsapp.InitApp(gaehost.GaeBotHost{})
}
