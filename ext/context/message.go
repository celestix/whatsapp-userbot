package context

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type Message struct {
	ctx context.Context
	*events.Message
}

func (m *Message) ArgsN(n int) []string {
	if m.Message.Message.Conversation != nil {
		return strings.SplitN(*m.Message.Message.Conversation, " ", n)
	} else if m.Message.Message.ExtendedTextMessage != nil {
		return strings.SplitN(*m.Message.Message.ExtendedTextMessage.Text, " ", n)
	}
	return []string{}
}

func (m *Message) Args() []string {
	if m.Message.Message.Conversation != nil {
		return strings.Fields(*m.Message.Message.Conversation)
	} else if m.Message.Message.ExtendedTextMessage != nil {
		return strings.Fields(*m.Message.Message.ExtendedTextMessage.Text)
	}
	return []string{}
}

func (m *Message) GetText() string {
	if m.Message.Message.Conversation != nil {
		return m.Message.Message.GetConversation()
	} else if m.Message.Message.ExtendedTextMessage != nil {
		return m.Message.Message.ExtendedTextMessage.GetText()
	}
	return ""
}

func (m *Message) Send(client *whatsmeow.Client, to types.JID, text string) (resp whatsmeow.SendResponse, err error) {
	return client.SendMessage(m.ctx, to, &proto.Message{
		Conversation: &text,
	})
}

func (m *Message) Reply(client *whatsmeow.Client, text string) (whatsmeow.SendResponse, error) {
	return client.SendMessage(m.ctx, m.Info.Chat, &proto.Message{
		ExtendedTextMessage: &proto.ExtendedTextMessage{
			Text: &text,
			ContextInfo: &proto.ContextInfo{
				StanzaId:      &m.Info.ID,
				Participant:   stringPtr(m.Info.Sender.String()),
				QuotedMessage: m.Message.Message,
			},
		},
	})
}

func (m *Message) Edit(client *whatsmeow.Client, text string) (whatsmeow.SendResponse, error) {
	return client.SendMessage(m.ctx, m.Info.Chat, client.BuildEdit(m.Info.Chat, m.Info.ID, &proto.Message{
		Conversation: &text,
	}))
}
