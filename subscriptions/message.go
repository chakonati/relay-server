package subscriptions

import (
	"log"
	"server/defs"
	"server/persistence"
	"sync"
)

type MessageSubscription struct {
	wg sync.WaitGroup
}

type MessageNotification struct {
	MessageID uint64 `key:"messageId"`
	From      string `key:"from"`
	DeviceID  int    `key:"deviceId"`
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
	s.NotifyUnconfirmed(subscriber.(MessageSubscriber))
}

func (s *MessageSubscription) Unsubscribe(subscriber Subscriber) {
	messageSubMut.Lock()
	defer messageSubMut.Unlock()
	delete(messageSubscribers, subscriber.(MessageSubscriber))
}

func (s *MessageSubscription) NotifyUnconfirmed(sub MessageSubscriber) {
	msgChan, errChan := persistence.Messages.StreamMessages()
	for msg := range msgChan {
		s.NotifyMessageSingleSubscriber(sub, msg)
	}
	for err := range errChan {
		log.Println(err)
	}
}

func (s *MessageSubscription) NotifyMessage(msg *defs.Message) {
	for _, subscriber := range messageSubscribers {
		s.wg.Add(1)
		go func(sub MessageSubscriber) {
			s.NotifyMessageSingleSubscriber(sub, msg)
			s.wg.Done()
		}(subscriber)
	}
	s.wg.Wait()
}

func (s *MessageSubscription) NotifyMessageSingleSubscriber(sub MessageSubscriber, msg *defs.Message) {
	sub.NotifyMessageNotification(s, &MessageNotification{
		MessageID: msg.ID,
		From:      msg.From,
		DeviceID:  msg.DeviceID,
	})
}

var _ Subscription = (*MessageSubscription)(nil)
var _ Subscriber = (*MessageSubscriber)(nil)
