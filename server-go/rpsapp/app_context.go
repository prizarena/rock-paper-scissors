package rpsapp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/strongo/app"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpsmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"reflect"
	"time"
	"github.com/strongo-games/rock-paper-scissors/server-go/rpstrans"
)

type rpsAppContext struct {
}

func (appCtx rpsAppContext) AppUserEntityKind() string {
	return rpsmodels.UserKind
}

func (appCtx rpsAppContext) AppUserEntityType() reflect.Type {
	return reflect.TypeOf(&rpsmodels.UserEntity{})
}

func (appCtx rpsAppContext) NewBotAppUserEntity() bots.BotAppUser {
	return &rpsmodels.UserEntity{
		DtCreated: time.Now(),
	}
}

func (appCtx rpsAppContext) NewAppUserEntity() strongo.AppUser {
	return appCtx.NewBotAppUserEntity()
}

func (appCtx rpsAppContext) GetTranslator(c context.Context) strongo.Translator {
	return strongo.NewMapTranslator(c, rpstrans.TRANS)
}

type LocalesProvider struct {
}

func (LocalesProvider) GetLocaleByCode5(code5 string) (strongo.Locale, error) {
	return strongo.LocaleEnUS, nil
}

func (appCtx rpsAppContext) SupportedLocales() strongo.LocalesProvider {
	return RpsLocalesProvider{}
}

type RpsLocalesProvider struct {
}

func (RpsLocalesProvider) GetLocaleByCode5(code5 string) (locale strongo.Locale, err error) {
	switch code5 {
	case strongo.LocaleCodeEnUS:
		return strongo.LocaleEnUS, nil
	case strongo.LocalCodeRuRu:
		return strongo.LocaleRuRu, nil
	default:
		return locale, errors.New("Unsupported locale: " + code5)
	}
}

var _ strongo.LocalesProvider = (*RpsLocalesProvider)(nil)

func (appCtx rpsAppContext) GetBotChatEntityFactory(platform string) func() bots.BotChat {
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

var _ bots.BotAppContext = (*rpsAppContext)(nil)

