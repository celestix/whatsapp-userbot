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

func addNote(client *whatsmeow.Client, ctx *context.Context) error {
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

func deleteNote(client *whatsmeow.Client, ctx *context.Context) error {
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

func getNote(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) == 1 {
		return ext.EndGroups
	}
	key := args[1]
	value := sql.GetNote(strings.ToLower(key)).Value
	if value == "" {
		value = "```Note not found```"
	}
	_, _ = ctx.Message.Edit(client, value)
	return ext.EndGroups
}

func getNoteHash(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) == 0 {
		return nil
	}
	text := strings.ToLower(args[0])
	if !strings.HasPrefix(text, "#") {
		return nil
	}
	text = text[1:]
	value := sql.GetNote(text).Value
	if value == "" {
		return nil
	}
	_, _ = ctx.Message.Edit(client, value)
	return nil
}

func listNotes(client *whatsmeow.Client, ctx *context.Context) error {
	text := "*List of notes*:"
	for _, note := range sql.GetNotes() {
		text += fmt.Sprintf("\n- ```%s```", note.Name)
	}
	if text == "*List of notes*:" {
		text = "No notes are present."
	}
	_, _ = ctx.Message.Edit(client, text)
	return ext.EndGroups
}

func (*Module) LoadNote(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("note")
	defer ppLogger.Println("Loaded Note module")
	dispatcher.AddHandler(
		handlers.NewCommand("add", authorizedOnly(addNote), ppLogger.Create("add-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Add a note."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("clear", authorizedOnly(deleteNote), ppLogger.Create("clear-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Clear a note."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("get", authorizedOnly(getNote), ppLogger.Create("get-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Get a note from key."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("notes", authorizedOnly(listNotes), ppLogger.Create("notes").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Display keys of all saved notes."),
	)
	dispatcher.AddHandlerToGroup(
		handlers.NewMessage(authorizedOnlyMessages(getNoteHash), ppLogger.Create("get-hash").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Get a note using #key format."),
		1,
	)
}
