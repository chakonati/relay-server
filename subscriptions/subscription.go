package subscriptions

import "log"

type Subscription interface {
	Subscribe(subscriber Subscriber)
	Unsubscribe(subscriber Subscriber)
}
type Subscriber interface{}

type SubscriptionList struct {
	MessagesSubscription *MessageSubscription
}

var Subscriptions SubscriptionList

func init() {
	log.Println("Initializing subscriptions")
	Subscriptions.MessagesSubscription = &MessageSubscription{}
	log.Println("Subscriptions successfully initialized")
}
