package rpsmodels

import "github.com/strongo/db"

type User struct {
	db.IntegerID
	*UserEntity
}

type UserEntity struct {

}
