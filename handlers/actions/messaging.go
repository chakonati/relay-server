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

func (m *MessageHandler) SendMessage(message []byte) error {
	return messaging.MessageReceived(&defs.Message{
		EncryptedMessage: message,
	})
}

func (m *MessageHandler) GetMessage(id uint64, password string) (*defs.Message, error) {
	if !m.setup.IsPasswordValid(password) {
		return nil, InvalidPasswordError()
	}

	return persistence.Messages.Message(id)
}

func (m *MessageHandler) NextMessageID(password string) (uint64, bool, error) {
	if !m.setup.IsPasswordValid(password) {
		return 0, false, InvalidPasswordError()
	}

	return persistence.Messages.FirstID()
}

func (m *MessageHandler) ConfirmReceived(id uint64, password string) error {
	if !m.setup.IsPasswordValid(password) {
		return InvalidPasswordError()
	}

	msg, err := persistence.Messages.Message(id)
	if err != nil {
		return errors.Wrapf(err, "could not get message information for ID %d", id)
	}
	if msg == nil {
		return fmt.Errorf("message with ID %d does not exist", id)
	}

	return persistence.Messages.RemoveMessageByID(id)
}
