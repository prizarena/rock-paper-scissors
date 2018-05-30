package rpsmodels

import (
	"testing"
	"github.com/strongo/bots-framework/core"
)

func TestUser(t *testing.T) {
	var _ bots.BotAppUser = (*User)(nil)
}
