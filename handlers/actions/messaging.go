package actions

import (
	"fmt"
	"server/defs"
	"server/messaging"
	"server/persistence"

	"github.com/pkg/errors"
)

type MessageHandler struct {
	setup SetupHandler
}

type MessagingSendMessageResponse struct{ Error error }

func (m *MessageHandler) SendMessage(message []byte) *MessagingSendMessageResponse {
	return &MessagingSendMessageResponse{messaging.MessageReceived(&defs.Message{
		EncryptedMessage: message,
	})}
}

type MessagingGetMessageResponse struct {
	Message *defs.Message
	Error   error
}

func (m *MessageHandler) GetMessage(id uint64, password string) *MessagingGetMessageResponse {
	if !m.setup.IsPasswordValid(password).Valid {
		return &MessagingGetMessageResponse{nil, InvalidPasswordError()}
	}

	msg, err := persistence.Messages.Message(id)
	return &MessagingGetMessageResponse{msg, err}
}

type MessagingNextMessageIdResponse struct {
	MessageId uint64
	Exists    bool
	Error     error
}

func (m *MessageHandler) NextMessageId(password string) *MessagingNextMessageIdResponse {
	if !m.setup.IsPasswordValid(password).Valid {
		return &MessagingNextMessageIdResponse{0, false, InvalidPasswordError()}
	}

	msgId, exists, err := persistence.Messages.FirstID()
	return &MessagingNextMessageIdResponse{msgId, exists, err}
}

type MessagingConfirmReceivedResponse struct{ Error error }

func (m *MessageHandler) ConfirmReceived(id uint64, password string) (response *MessagingConfirmReceivedResponse) {
	response = &MessagingConfirmReceivedResponse{}

	if !m.setup.IsPasswordValid(password).Valid {
		response.Error = InvalidPasswordError()
		return
	}

	msg, err := persistence.Messages.Message(id)
	if err != nil {
		response.Error = errors.Wrapf(err, "could not get message information for ID %d", id)
		return
	}
	if msg == nil {
		response.Error = fmt.Errorf("message with ID %d does not exist", id)
		return
	}

	response.Error = persistence.Messages.RemoveMessageByID(id)
	return
}
