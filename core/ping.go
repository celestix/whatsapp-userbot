package core

import (
	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	waLogger "github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow"
)

func ping(client *whatsmeow.Client, ctx *context.Context) error {
	_, _ = ctx.Message.Edit(client, "Pong!")
	return ext.EndGroups
}

func (*Module) LoadPing(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("ping")
	defer ppLogger.Println("Loaded Ping module")
	dispatcher.AddHandler(
		handlers.NewCommand("ping", authorizedOnly(ping), ppLogger.Create("ping-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("To ping the userbot."),
	)
}
