package subscriptions

import (
	"server/defs"
	"sync"
)

type MessageSubscription struct {
	wg sync.WaitGroup
}

type MessageNotification struct {
	MessageID uint64
}

type MessageSubscriber interface {
	NotifyMessageNotification(sub *MessageSubscription, notification *MessageNotification)
}

var messageSubMut sync.Mutex
var messageSubscribers = map[MessageSubscriber]MessageSubscriber{}

func (s *MessageSubscription) Subscribe(subscriber Subscriber) {
	messageSubMut.Lock()
	defer messageSubMut.Unlock()
	messageSubscribers[subscriber.(MessageSubscriber)] = subscriber.(MessageSubscriber)
}

func (s *MessageSubscription) Unsubscribe(subscriber Subscriber) {
	messageSubMut.Lock()
	defer messageSubMut.Unlock()
	delete(messageSubscribers, subscriber.(MessageSubscriber))
}

func (s *MessageSubscription) NotifyMessage(msg *defs.Message) {
	for _, subscriber := range messageSubscribers {
		s.wg.Add(1)
		go func(sub MessageSubscriber) {
			sub.NotifyMessageNotification(s, &MessageNotification{
				MessageID: msg.ID,
			})
			s.wg.Done()
		}(subscriber)
	}
	s.wg.Wait()
}

var _ Subscription = (*MessageSubscription)(nil)
var _ Subscriber = (*MessageSubscriber)(nil)
