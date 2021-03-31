package persistence

import (
	"log"
	"server/defs"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type MessageDAO struct{}

const messageBucketName = "messages"

var Messages = MessageDAO{}

func (m *MessageDAO) init() error {
	if err := initBucket(&messages, messageBucketName); err != nil {
		return err
	}
	return nil
}

func (m *MessageDAO) bucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(messageBucketName))
}

func (m *MessageDAO) AddMessage(message *defs.Message) error {
	return messages.db.Update(func(tx *bbolt.Tx) error {
		var err error
		b := m.bucket(tx)
		message.ID, err = b.NextSequence()
		if err != nil {
			return errors.Wrap(err, "could not get next seq for message")
		}
		if err := PutStruct(b, UInt64Bytes(message.ID), message); err != nil {
			return errors.Wrap(err, "could not store pre-key bundle")
		}
		log.Println("The pre-key bundle is now registered.")
		return nil
	})
}

func (m *MessageDAO) Message(ID uint64) (*defs.Message, error) {
	var msg *defs.Message
	err := messages.db.View(func(tx *bbolt.Tx) error {
		return GetStruct(m.bucket(tx), UInt64Bytes(ID), msg)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not get message for ID %d", ID)
	}
	return msg, nil
}

func (m *MessageDAO) RemoveMessageByID(id uint64) error {
	return messages.db.Update(func(tx *bbolt.Tx) error {
		return m.bucket(tx).Delete(UInt64Bytes(id))
	})
}

func (m *MessageDAO) RemoveMessage(message *defs.Message) error {
	return m.RemoveMessageByID(message.ID)
}

func (m *MessageDAO) FirstID() (uint64, bool, error) {
	var id uint64
	var hasMessages bool
	err := messages.db.View(func(tx *bbolt.Tx) error {
		keyByt, _ := m.bucket(tx).Cursor().First()
		hasMessages = keyByt != nil
		if hasMessages {
			id = BytesUInt64(keyByt)
		}
		return nil
	})
	if err != nil {
		return 0, false, errors.Wrap(err, "could not get first ID")
	}
	return id, hasMessages, nil
}
