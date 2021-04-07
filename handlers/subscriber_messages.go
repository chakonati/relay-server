package handlers

import (
	"log"
	"server/subscriptions"
)

func (h *Handler) NotifyMessageNotification(
	sub *subscriptions.MessageSubscription,
	notification *subscriptions.MessageNotification,
) {
	if err := h.sendMessage(&Notification{
		MessageType: MessageTypeOneway,
		Data:        notification,
	}); err != nil {
		log.Println("Note: could not notify of message:", err)
	}
}

var _ subscriptions.MessageSubscriber = (*Handler)(nil)
