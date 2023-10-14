package context

import (
	"context"

	"github.com/celestix/whatsapp-userbot/logger"

	"go.mau.fi/whatsmeow/types/events"
)

func stringPtr(s string) *string {
	return &s
}

type Context struct {
	Message *Message
	Logger  *logger.Logger
}

func New(ctx context.Context, evt interface{}) *Context {
	switch v := evt.(type) {
	case *events.Message:
		return &Context{
			Message: &Message{
				Message: v,
				ctx:     ctx,
			},
		}
	}
	return &Context{}
}
