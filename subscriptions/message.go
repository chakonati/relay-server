package subscriptions

import "sync"

type MessageSubscription struct{}

type MessageNotification struct{}

type MessageSubscriber interface {
	NotifyMessageNotification(sub MessageSubscription, notification *MessageNotification)
}

var messageSubMut sync.Mutex
var messageSubscribers = map[MessageSubscriber]MessageSubscriber{}

func (s MessageSubscription) Subscribe(subscriber Subscriber) {
	messageSubMut.Lock()
	defer messageSubMut.Unlock()
	messageSubscribers[subscriber.(MessageSubscriber)] = subscriber.(MessageSubscriber)
}

func (s MessageSubscription) Unsubscribe(subscriber Subscriber) {
	messageSubMut.Lock()
	defer messageSubMut.Unlock()
	delete(messageSubscribers, subscriber.(MessageSubscriber))
}

var _ Subscription = (*MessageSubscription)(nil)
var _ Subscriber = (*MessageSubscriber)(nil)
