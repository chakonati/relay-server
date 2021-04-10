package messaging

import (
	"server/defs"
	"server/persistence"
	"time"

	"github.com/pkg/errors"
)

var InboundQueue = make(chan *defs.Message)
var OutboundQueue = make(chan *defs.Message)

func MessageReceived(msg *defs.Message) error {
	if len(msg.EncryptedMessage) == 0 {
		return errors.New("Your message is empty. Make sure you actually put data in it. Rejected.")
	}

	msg.At = time.Now()
	if err := persistence.Messages.AddMessage(msg); err != nil {
		return errors.Wrap(err, "could not relay message")
	}
	InboundQueue <- msg
	return nil
}
