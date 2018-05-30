package rpsmodels

import (
	"github.com/strongo/db"
	"github.com/strongo/app"
	"time"
	"github.com/strongo/bots-framework/core"
)

const UserKind = "User"

type User struct {
	db.IntegerID
	*UserEntity
}

type UserEntity struct {
	strongo.AppUserBase
	DtCreated time.Time
}


var _ bots.BotAppUser = (*User)(nil)