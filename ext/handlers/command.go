package handlers

import (
	"fmt"
	"strings"

	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow"
)

var DefaultTriggers = []rune{'.'}

type Command struct {
	Prefix      []rune
	Command     string
	Logger      *logger.Logger
	Response    Response
	description string
}

func NewCommand(command string, callback Response, logger *logger.Logger) *Command {
	return &Command{
		Prefix:   DefaultTriggers,
		Command:  strings.ToLower(command),
		Logger:   logger,
		Response: callback,
	}
}

func NewCommandWithPrefix(prefix []rune, command string, callback Response, logger *logger.Logger) *Command {
	return &Command{
		Prefix:   prefix,
		Command:  strings.ToLower(command),
		Logger:   logger,
		Response: callback,
	}
}

func (m *Command) AddDescription(desc string) Handler {
	m.description = fmt.Sprintf("```.%s```: %s", m.Command, desc)
	return m
}

func (m Command) GetDescription() string {
	return m.description
}

func (m Command) GetName() string {
	return m.Command
}

func (m Command) CheckUpdate(ctx *context.Context) bool {
	return ctx.Message != nil
}

func (m Command) HandleUpdate(client *whatsmeow.Client, ctx *context.Context) error {
	msg := ctx.Message.Message
	var text string
	if msg.Message.Conversation != nil {
		text = strings.Fields(*msg.Message.Conversation)[0]
	} else if msg.Message.ExtendedTextMessage != nil {
		text = strings.Fields(*msg.Message.ExtendedTextMessage.Text)[0]
	} else {
		return nil
	}

	for _, prefix := range m.Prefix {
		if prefix == rune(text[0]) && m.Command == strings.ToLower(text[1:]) {
			ctx.Logger = m.Logger
			return m.Response(client, ctx)
		}
	}
	return nil
}
