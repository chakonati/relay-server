package messaging

import "server/subscriptions"

type Relay struct{}

func (r *Relay) StartWorking() {
	for msg := range InboundQueue {
		subscriptions.Subscriptions.MessagesSubscription.NotifyMessage(msg)
	}
}
