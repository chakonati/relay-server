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
	EncryptedMessage []byte
}

func (m *MessageHandler) SendMessage(request MessagingSendMessageRequest) *MessagingSendMessageResponse {
	return &MessagingSendMessageResponse{messaging.MessageReceived(&defs.Message{
		EncryptedMessage: request.EncryptedMessage,
	})}
}

type MessagingGetMessageResponse struct {
	MessageId        uint64
	EncryptedMessage []byte
	Error            error
}

type MessagingGetMessageRequest struct {
	MessageId uint64
	Password  string
}

func (m *MessageHandler) GetMessage(request MessagingGetMessageRequest) *MessagingGetMessageResponse {
	if !m.setup.isPasswordValid(request.Password).Valid {
		return &MessagingGetMessageResponse{Error: InvalidPasswordError()}
	}

	msg, err := persistence.Messages.Message(request.MessageId)
	if err != nil {
		return &MessagingGetMessageResponse{Error: err}
	}
	if msg.EncryptedMessage == nil {
		return &MessagingGetMessageResponse{Error: errors.New("encrypted message is nil, must not be nil")}
	}

	return &MessagingGetMessageResponse{
		MessageId:        msg.ID,
		EncryptedMessage: msg.EncryptedMessage,
	}
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
	MessageId uint64
	Password  string
}

func (m *MessageHandler) ConfirmReceived(request MessagingConfirmReceivedRequest) (response *MessagingConfirmReceivedResponse) {
	response = &MessagingConfirmReceivedResponse{}

	if !m.setup.isPasswordValid(request.Password).Valid {
		response.Error = InvalidPasswordError()
		return
	}

	msg, err := persistence.Messages.Message(request.MessageId)
	if err != nil {
		response.Error = errors.Wrapf(err, "could not get Message information for ID %d", request.MessageId)
		return
	}
	if msg == nil {
		response.Error = fmt.Errorf("message with ID %d does not exist", request.MessageId)
		return
	}

	response.Error = persistence.Messages.RemoveMessageByID(request.MessageId)
	return
}
