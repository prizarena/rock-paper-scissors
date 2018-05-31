package rpsmodels

import "github.com/strongo/db"

const GameCanvasKind = "GameCanvas"

type GameCanvasEntity struct {
	UserIDs   []string
	UserMoves []string
}

type GameCanvas struct {
	db.StringID
	*GameCanvasEntity
}

var _ db.EntityHolder = (*GameCanvas)(nil)

func (GameCanvas) Kind() string {
	return GameCanvasKind
}

func (canvas *GameCanvas) SetEntity(v interface{}) {
	canvas.GameCanvasEntity = v.(*GameCanvasEntity)
}

func (canvas GameCanvas) Entity() interface{} {
	return canvas.GameCanvasEntity
}

func (canvas GameCanvas) NewEntity() interface{} {
	return &GameCanvasEntity{}
}
