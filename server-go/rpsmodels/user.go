package rpsmodels

import (
	"github.com/strongo/db"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/app"
)

const UserKind = "User"

type User struct {
	db.IntegerID
	*UserEntity
}

type UserEntity = strongo.AppUserBase

var _ bots.BotAppUser = (*User)(nil)