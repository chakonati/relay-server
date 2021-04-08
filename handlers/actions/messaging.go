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

type MessagingSendMessageRequest struct {
	Message []byte
}

func (m *MessageHandler) SendMessage(request MessagingSendMessageRequest) *MessagingSendMessageResponse {
	return &MessagingSendMessageResponse{messaging.MessageReceived(&defs.Message{
		EncryptedMessage: request.Message,
	})}
}

type MessagingGetMessageResponse struct {
	Message *defs.Message
	Error   error
}

type MessagingGetMessageRequest struct {
	Id       uint64
	Password string
}

func (m *MessageHandler) GetMessage(request MessagingGetMessageRequest) *MessagingGetMessageResponse {
	if !m.setup.isPasswordValid(request.Password).Valid {
		return &MessagingGetMessageResponse{nil, InvalidPasswordError()}
	}

	msg, err := persistence.Messages.Message(request.Id)
	return &MessagingGetMessageResponse{msg, err}
}

type MessagingNextMessageIdResponse struct {
	MessageId uint64
	Exists    bool
	Error     error
}

type MessagingNextMessageIdRequest struct {
	Password string
}

func (m *MessageHandler) NextMessageId(request MessagingNextMessageIdRequest) *MessagingNextMessageIdResponse {
	if !m.setup.isPasswordValid(request.Password).Valid {
		return &MessagingNextMessageIdResponse{0, false, InvalidPasswordError()}
	}

	msgId, exists, err := persistence.Messages.FirstID()
	return &MessagingNextMessageIdResponse{msgId, exists, err}
}

type MessagingConfirmReceivedResponse struct{ Error error }

type MessagingConfirmReceivedRequest struct {
	Id       uint64
	Password string
}

func (m *MessageHandler) ConfirmReceived(request MessagingConfirmReceivedRequest) (response *MessagingConfirmReceivedResponse) {
	response = &MessagingConfirmReceivedResponse{}

	if !m.setup.isPasswordValid(request.Password).Valid {
		response.Error = InvalidPasswordError()
		return
	}

	msg, err := persistence.Messages.Message(request.Id)
	if err != nil {
		response.Error = errors.Wrapf(err, "could not get Message information for ID %d", request.Id)
		return
	}
	if msg == nil {
		response.Error = fmt.Errorf("Message with ID %d does not exist", request.Id)
		return
	}

	response.Error = persistence.Messages.RemoveMessageByID(request.Id)
	return
}
