package subscriptions

type Subscription interface {
	Subscribe(subscriber Subscriber)
	Unsubscribe(subscriber Subscriber)
}
type Subscriber interface{}

type SubscriptionList struct {
	MessagesSubscription MessageSubscription
}

var Subscriptions SubscriptionList
