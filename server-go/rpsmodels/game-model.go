package rpsmodels

import (
	"github.com/strongo/db"
	"github.com/prizarena/arena/canvas"
)

const RpsGameKind = "G"

type RpsGame struct {
	db.IntegerID
	*RpsGameEntity
}

var _ db.EntityHolder = (*RpsGame)(nil)

type RpsGameEntity struct {
	BoardID string `datastore:",noindex,omitempty"`
	canvas.BoardEntity
}

func (RpsGame) Kind() string {
	return RpsGameKind
}

func (g RpsGame) NewEntity() interface{} {
	return &RpsGameEntity{}
}

func (g RpsGame) Entity() interface{} {
	return g.RpsGameEntity
}

func (g *RpsGame) SetEntity(v interface{}) {
	g.RpsGameEntity = v.(*RpsGameEntity)
}
