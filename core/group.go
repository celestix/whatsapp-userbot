package core

import (
	"fmt"

	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	waLogger "github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func info(client *whatsmeow.Client, ctx *context.Context) error {
	jid := ctx.Message.Info.Chat
	info, err := client.GetGroupInfo(jid)
	if err != nil {
		_, _ = ctx.Message.Edit(client, err.Error())
		return ext.EndGroups
	}
	if info.Topic == "" {
		info.Topic = "None"
	}
	text := "*Group Info*"
	text += fmt.Sprintf("\n*Name*: ```%s```", info.Name)
	text += fmt.Sprintf("\n*Owner*: %s", info.OwnerJID.User)
	text += fmt.Sprintf("\n*Topic*: ```%s```", info.Topic)
	text += fmt.Sprintf("\n*ID*: ```%s```", jid.String())
	_, _ = ctx.Message.Edit(client, text)
	return ext.EndGroups
}

func linkinfo(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) <= 1 {
		return ext.EndGroups
	}
	info, err := client.GetGroupInfoFromLink(args[1])
	if err != nil {
		_, _ = ctx.Message.Edit(client, err.Error())
		return ext.EndGroups
	}
	if info.Topic == "" {
		info.Topic = "None"
	}
	text := "*Group Info*"
	text += fmt.Sprintf("\n*Name*: ```%s```", info.Name)
	text += fmt.Sprintf("\n*Owner*: %s", info.OwnerJID.User)
	text += fmt.Sprintf("\n*Topic*: ```%s```", info.Topic)
	text += fmt.Sprintf("\n*ID*: ```%s```", info.JID.String())
	_, _ = ctx.Message.Edit(client, text)
	return ext.EndGroups
}

func addUser(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.Args()
	if len(args) <= 1 {
		return ext.EndGroups
	}
	if ctx.Message.Info.IsGroup {
		return ext.EndGroups
	}
	jid, err := types.ParseJID(args[1])
	if err != nil {
		_, _ = ctx.Message.Edit(client, "failed to parse JID: "+err.Error())
		return ext.EndGroups
	}
	_, err = client.UpdateGroupParticipants(jid, []types.JID{ctx.Message.Info.Chat}, whatsmeow.ParticipantChangeAdd)
	if err != nil {
		_, _ = ctx.Message.Edit(client, "failed to add: "+err.Error())
		return ext.EndGroups
	}
	_, _ = ctx.Message.Edit(client, "Added...")
	return ext.EndGroups
}

func removeUser(client *whatsmeow.Client, ctx *context.Context) error {
	msg := ctx.Message.Message
	if !msg.Info.IsGroup {
		return ext.EndGroups
	}
	var jidstring string
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		jidstring = *msg.Message.ExtendedTextMessage.ContextInfo.Participant

	} else if len(ctx.Message.Args()) > 1 {
		jidstring = ctx.Message.Args()[1]
	}
	jid, err := types.ParseJID(jidstring)
	if err != nil {
		fmt.Println(err.Error())
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("failed to get user: %s", err.Error()))
		return ext.EndGroups
	}
	_, err = client.UpdateGroupParticipants(jid, []types.JID{ctx.Message.Info.Chat}, whatsmeow.ParticipantChangeRemove)
	if err != nil {
		_, _ = ctx.Message.Edit(client, "failed to remove: "+err.Error())
		return ext.EndGroups
	}
	_, _ = ctx.Message.Edit(client, "Removed...")

	return ext.EndGroups
}

func promoteUser(client *whatsmeow.Client, ctx *context.Context) error {
	msg := ctx.Message.Message
	if !msg.Info.IsGroup {
		return ext.EndGroups
	}
	var jidstring string
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		jidstring = *msg.Message.ExtendedTextMessage.ContextInfo.Participant

	} else if len(ctx.Message.Args()) > 1 {
		jidstring = ctx.Message.Args()[1]
	}
	jid, err := types.ParseJID(jidstring)
	if err != nil {
		fmt.Println(err.Error())
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("failed to get user: %s", err.Error()))
		return ext.EndGroups
	}
	_, err = client.UpdateGroupParticipants(jid, []types.JID{ctx.Message.Info.Chat}, whatsmeow.ParticipantChangePromote)
	if err != nil {
		_, _ = ctx.Message.Edit(client, "failed to promote: "+err.Error())
		return ext.EndGroups
	}
	_, _ = ctx.Message.Edit(client, "Promoted...")

	return ext.EndGroups
}

func demoteUser(client *whatsmeow.Client, ctx *context.Context) error {
	msg := ctx.Message.Message
	if !msg.Info.IsGroup {
		return ext.EndGroups
	}
	var jidstring string
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		jidstring = *msg.Message.ExtendedTextMessage.ContextInfo.Participant

	} else if len(ctx.Message.Args()) > 1 {
		jidstring = ctx.Message.Args()[1]
	}
	jid, err := types.ParseJID(jidstring)
	if err != nil {
		fmt.Println(err.Error())
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("failed to get user: %s", err.Error()))
		return ext.EndGroups
	}
	_, err = client.UpdateGroupParticipants(jid, []types.JID{ctx.Message.Info.Chat}, whatsmeow.ParticipantChangePromote)
	if err != nil {
		_, _ = ctx.Message.Edit(client, "failed to demote: "+err.Error())
		return ext.EndGroups
	}
	_, _ = ctx.Message.Edit(client, "Demoted...")

	return ext.EndGroups
}

func (*Module) LoadGroups(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("groups")
	defer ppLogger.Println("Loaded Groups module")
	dispatcher.AddHandler(
		handlers.NewCommand("info", authorizedOnly(info), ppLogger.Create("info-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription(`Get information about a group.`),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("linkinfo", authorizedOnly(linkinfo), ppLogger.Create("linkinfo-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription(`Get information about a group using chat link.`),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("adduser", authorizedOnly(addUser), ppLogger.Create("adduser-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Add user to a group."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("rmuser", authorizedOnly(removeUser), ppLogger.Create("rmuser-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Remove user from a group."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("promote", authorizedOnly(promoteUser), ppLogger.Create("rmuser-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Promote a civilian to an admin ."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("demote", authorizedOnly(demoteUser), ppLogger.Create("rmuser-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Demote an admin to a civilian."),
	)
}
