package core

import (
	"reflect"

	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	"github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow"
)

var LOGGER = logger.NewLogger(logger.LevelInfo)

type Module struct{}

func stringPtr(s string) *string {
	return &s
}

func authorizedOnly(callback handlers.Response) handlers.Response {
	return func(client *whatsmeow.Client, ctx *context.Context) error {
		if !ctx.Message.Info.IsFromMe {
			return ext.EndGroups
		}
		return callback(client, ctx)
	}
}

func authorizedOnlyMessages(callback handlers.Response) handlers.Response {
	return func(client *whatsmeow.Client, ctx *context.Context) error {
		if !ctx.Message.Info.IsFromMe {
			return nil
		}
		return callback(client, ctx)
	}
}

func Load(dispatcher *ext.Dispatcher) {
	defer LOGGER.Println("Loaded all modules")
	Type := reflect.TypeOf(&Module{})
	Value := reflect.ValueOf(&Module{})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(dispatcher)})
	}
}
