package messaging

import (
	"log"
	"server/subscriptions"
)

type Relay struct{}

func (r *Relay) StartWorking() {
	log.Println("Inbound relay is working")
	for msg := range InboundQueue {
		subscriptions.Subscriptions.MessagesSubscription.NotifyMessage(msg)
	}
}
