package handlers

import (
	"log"
	"server/encoders"
	"server/subscriptions"
)

func (h *Handler) NotifyMessageNotification(
	sub *subscriptions.MessageSubscription,
	notification *subscriptions.MessageNotification,
) {
	log.Println("Notifying of message notification, message ID:", notification.MessageID)
	byt, err := encoders.Msgpack{}.Marshal(notification)
	if err != nil {
		log.Println("Warning: could not encode message notification")
		return
	}
	if err := h.sendMessage(&Notification{
		MessageType:      MessageTypeNotification,
		SubscriptionName: "messages",
		Data:             byt,
	}); err != nil {
		log.Println("Note: could not notify of message:", err)
	}
}

var _ subscriptions.MessageSubscriber = (*Handler)(nil)
