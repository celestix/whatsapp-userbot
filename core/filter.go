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

func addFilter(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) <= 1 {
		return ext.EndGroups
	}
	msg := ctx.Message.Message
	key := args[1]
	var value string
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		qmsg := msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
		switch {
		case qmsg.Conversation != nil:
			value = qmsg.GetConversation()
		case qmsg.ExtendedTextMessage != nil:
			value = qmsg.ExtendedTextMessage.GetText()
		}
	} else {
		if len(args) == 2 {
			return ext.EndGroups
		}
		value = strings.Join(args[2:], " ")
	}
	go sql.AddNote(strings.ToLower(key), value)
	_, _ = ctx.Message.Edit(client, fmt.Sprintf("Added Note ```%s```.", key))
	return ext.EndGroups
}

func removeFilter(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) == 1 {
		return ext.EndGroups
	}
	key := args[1]
	if sql.DeleteNote(strings.ToLower(key)) {
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("Successfully delete note '```%s```'.", key))
	} else {
		_, _ = ctx.Message.Edit(client, "Failed to delete that note!")
	}
	return ext.EndGroups
}

func loadFilters(l *waLogger.Logger, dispatcherGroup int) {

}

func listFilters(client *whatsmeow.Client, ctx *context.Context) error {
	text := "*List of filters*:"
	for _, note := range sql.GetNotes() {
		text += fmt.Sprintf("\n- ```%s```", note.Name)
	}
	if text == "*List of filters*:" {
		text = "No notes are present."
	}
	_, _ = ctx.Message.Edit(client, text)
	return ext.EndGroups
}

func (*Module) LoadFilter(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("filter")
	defer ppLogger.Println("Loaded Note module")
	dispatcher.AddHandler(
		handlers.NewCommand("filter", authorizedOnly(addFilter), ppLogger.Create("filter-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Add a filter."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("remove", authorizedOnly(removeFilter), ppLogger.Create("remove-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Remove a filter."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("filters", authorizedOnly(listFilters), ppLogger.Create("list-filters").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Display keys of all saved notes."),
	)
	loadFilters(
		ppLogger.Create("filter-cmd").
			ChangeLevel(waLogger.LevelInfo),
		1,
	)
}
