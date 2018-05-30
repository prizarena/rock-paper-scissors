package rpsapp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/strongo/app"
	// "github.com/strongo-games/bidding-tictactoe/server-go/rsp-trans"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"reflect"
	"time"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpstrans"
)

type BiddingTicTacToeAppContext struct {
}

func (appCtx BiddingTicTacToeAppContext) AppUserEntityKind() string {
	return rpsmodels.UserKind
}

func (appCtx BiddingTicTacToeAppContext) AppUserEntityType() reflect.Type {
	return reflect.TypeOf(&rpsmodels.UserEntity{})
}

func (appCtx BiddingTicTacToeAppContext) NewBotAppUserEntity() bots.BotAppUser {
	return &rpsmodels.UserEntity{
		DtCreated: time.Now(),
	}
}

func (appCtx BiddingTicTacToeAppContext) NewAppUserEntity() strongo.AppUser {
	return appCtx.NewBotAppUserEntity()
}

func (appCtx BiddingTicTacToeAppContext) GetTranslator(c context.Context) strongo.Translator {
	return strongo.NewMapTranslator(c, rpstrans.TRANS)
}

type LocalesProvider struct {
}

func (LocalesProvider) GetLocaleByCode5(code5 string) (strongo.Locale, error) {
	return strongo.LocaleEnUS, nil
}

func (appCtx BiddingTicTacToeAppContext) SupportedLocales() strongo.LocalesProvider {
	return BtttLocalesProvider{}
}

type BtttLocalesProvider struct {
}

func (BtttLocalesProvider) GetLocaleByCode5(code5 string) (locale strongo.Locale, err error) {
	switch code5 {
	case strongo.LocaleCodeEnUS:
		return strongo.LocaleEnUS, nil
	case strongo.LocalCodeRuRu:
		return strongo.LocaleRuRu, nil
	default:
		return locale, errors.New("Unsupported locale: " + code5)
	}
}

var _ strongo.LocalesProvider = (*BtttLocalesProvider)(nil)

func (appCtx BiddingTicTacToeAppContext) GetBotChatEntityFactory(platform string) func() bots.BotChat {
	switch platform {
	case telegram.PlatformID:
		return func() bots.BotChat {
			return &telegram.ChatEntity{
				TgChatEntityBase: *telegram.NewTelegramChatEntity(),
			}
		}
	default:
		panic("Unknown platform: " + platform)
	}
}

var _ bots.BotAppContext = (*BiddingTicTacToeAppContext)(nil)

var TheAppContext = BiddingTicTacToeAppContext{}
