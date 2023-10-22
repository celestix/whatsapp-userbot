package core

import (
	"fmt"
	"strings"

	"github.com/celestix/whatsapp-userbot/core/sql"
	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	waLogger "github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow"
)

func afk(client *whatsmeow.Client, ctx *context.Context) error {
	var reason string
	if len(ctx.Message.Args()) > 1 {
		reason = strings.Join(ctx.Message.Args()[1:], " ")
	}
	go sql.ToggleAfk(true, reason)
	_, _ = ctx.Message.Edit(client, "*Successfully enabled AFK mode.*")
	return ext.EndGroups
}

func workerAFK(client *whatsmeow.Client, ctx *context.Context) error {
	if ctx.Message.Info.IsGroup {
		return nil
	}
	status := sql.GetAfkStatus()
	if !ctx.Message.Info.IsFromMe {
		if !status.Working {
			return nil
		}
		text := "*I am currently Away-From-Keyboard.*"
		if status.Reason != "" {
			text += fmt.Sprintf("\n*Reason*: ```%s```", status.Reason)
		}
		_, _ = ctx.Message.Reply(client, text)
		return nil
	}
	if status.Working {
		sql.ToggleAfk(false, "")
		ctx.Logger.Println("Disabled AFK")
	}
	return nil
}

func silentAFK(client *whatsmeow.Client, ctx *context.Context) error {
	chatId := ctx.Message.Info.Chat.String()
	args := ctx.Message.Args()
	var silentAFK bool
	var text = "Turned off AFK mode for this chat."
	if len(args) == 1 {
		silentAFK = true
	} else {
		switch strings.ToLower(args[1]) {
		case "true", "yes", "on":
			silentAFK = true
		case "false", "no", "off":
			silentAFK = false
			text = "Turned on AFK mode for this chat."
		default:
			return ext.EndGroups
		}
	}
	sql.ShouldSilentAFK(chatId, silentAFK)
	ctx.Message.Edit(client, text)
	return ext.EndGroups
}

func (*Module) LoadAfk(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("afk")
	defer ppLogger.Println("Loaded Afk module")
	dispatcher.AddHandler(
		handlers.NewCommand("afk", authorizedOnly(afk), ppLogger.Create("afk-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription(`Turn on/off afk mode

AFK stands for Away From Keyboard, which should be used when you're offline for some work. 
		`),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("silentafk", authorizedOnly(silentAFK), ppLogger.Create("silentafk-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription(`Puts AFK on silent mode for a chat.`),
	)
	dispatcher.AddHandlerToGroup(
		handlers.NewMessage(workerAFK, ppLogger.Create("afk-worker").
			ChangeLevel(waLogger.LevelInfo),
		),
		2,
	)
}
