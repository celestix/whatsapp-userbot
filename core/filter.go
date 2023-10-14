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

func addFilter(dispatcher *ext.Dispatcher, dispatcherGroup int) func(
	client *whatsmeow.Client, ctx *context.Context,
) error {
	return func(client *whatsmeow.Client, ctx *context.Context) error {
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
		go func() {
			sql.AddFilter(strings.ToLower(key), value)
			initFilter(ctx.Logger, &sql.Filter{Name: key, Value: value}, dispatcher, dispatcherGroup)
		}()
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("Added Filter ```%s```.", key))
		return ext.EndGroups
	}
}

func removeFilter(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) == 1 {
		return ext.EndGroups
	}
	key := args[1]
	if sql.DeleteFilter(strings.ToLower(key)) {
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("Successfully delete filter '```%s```'.", key))
	} else {
		_, _ = ctx.Message.Edit(client, "Failed to delete that filter!")
	}
	return ext.EndGroups
}

func initFilter(l *waLogger.Logger, filter *sql.Filter, dispatcher *ext.Dispatcher, dispatcherGroup int) {
	dispatcher.AddHandlerToGroup(
		handlers.NewMessage(
			func(client *whatsmeow.Client, ctx *context.Context) error {
				text := ctx.Message.GetText()
				if !strings.Contains(strings.ToLower(text), filter.Name) {
					return ext.EndGroups
				}
				_, _ = ctx.Message.Reply(client, filter.Value)
				return ext.EndGroups
			},
			l.Create("filter-"+filter.Name).ChangeLevel(waLogger.LevelInfo),
		),
		dispatcherGroup,
	)
}

func loadFilters(l *waLogger.Logger, dispatcher *ext.Dispatcher, dispatcherGroup int) {
	for _, filter := range sql.GetFilters() {
		initFilter(l, &filter, dispatcher, dispatcherGroup)
	}
}

func listFilters(client *whatsmeow.Client, ctx *context.Context) error {
	text := "*List of filters*:"
	for _, filter := range sql.GetFilters() {
		text += fmt.Sprintf("\n- ```%s```", filter.Name)
	}
	if text == "*List of filters*:" {
		text = "No filters are present."
	}
	_, _ = ctx.Message.Edit(client, text)
	return ext.EndGroups
}

func (*Module) LoadFilter(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("filter")
	defer ppLogger.Println("Loaded Filter module")
	dispatcher.AddHandler(
		handlers.NewCommand("filter", authorizedOnly(addFilter(dispatcher, 1)), ppLogger.Create("filter-cmd").
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
		).AddDescription("Display keys of all saved filters."),
	)
	loadFilters(
		ppLogger.Create("filter-cmd").
			ChangeLevel(waLogger.LevelInfo),
		dispatcher,
		1,
	)
	// take you back
}
