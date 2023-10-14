package core

import (
	"fmt"

	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	waLogger "github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow"
)

func block(client *whatsmeow.Client, ctx *context.Context) error {
	if ctx.Message.Info.IsGroup {
		return ext.EndGroups
	}
	chat := ctx.Message.Info.Chat
	if _, err := client.UpdateBlocklist(chat, "block"); err != nil {
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("failed to block: %s", err.Error()))
	}
	return ext.EndGroups
}

func (*Module) LoadUserActions(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("user_actions")
	defer ppLogger.Println("Loaded UserActions module")
	dispatcher.AddHandler(
		handlers.NewCommand("block", authorizedOnly(block), ppLogger.Create("block-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Block a user."),
	)
}
