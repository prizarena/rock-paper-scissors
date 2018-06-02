package rpsgaedal

import (
	"github.com/strongo/db/gaedb"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsdal"
)

func RegisterDal() {
	rpsdal.DB = gaedb.NewDatabase()
}
