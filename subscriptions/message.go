package subscriptions

type MessageSubscription struct{}

var _ Subscription = (*MessageSubscription)(nil)
