package handlers

import (
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/logger"
	"github.com/celestix/whatsapp-userbot/utils"

	"go.mau.fi/whatsmeow"
)

type Message struct {
	Logger      *logger.Logger
	Response    Response
	description string
}

func NewMessage(callback Response, logger *logger.Logger) *Message {
	return &Message{
		Logger:   logger,
		Response: callback,
	}
}

func (m *Message) AddDescription(desc string) Handler {
	m.description = desc
	return m
}

func (m Message) GetDescription() string {
	return m.description
}

func (m Message) GetName() string {
	return utils.GetFuncName(m.Response)
}

func (m Message) CheckUpdate(ctx *context.Context) bool {
	return ctx.Message != nil
}

func (m Message) HandleUpdate(client *whatsmeow.Client, ctx *context.Context) error {
	ctx.Logger = m.Logger
	return m.Response(client, ctx)
}
