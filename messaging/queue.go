package messaging

import (
	"server/defs"
	"server/persistence"
	"time"

	"github.com/pkg/errors"
)

var InboundQueue chan *defs.Message
var OutboundQueue chan *defs.Message

func MessageReceived(msg *defs.Message) error {
	msg.At = time.Now()
	if err := persistence.Messages.AddMessage(msg); err != nil {
		return errors.Wrap(err, "could not relay message")
	}
	InboundQueue <- msg
	return nil
}
