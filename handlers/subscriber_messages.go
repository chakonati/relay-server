package handlers

import "server/subscriptions"

func (h *Handler) NotifyMessageNotification(
	sub subscriptions.MessageSubscription,
	notification *subscriptions.MessageNotification,
) {
	panic("implement me")
}

var _ subscriptions.MessageSubscriber = (*Handler)(nil)
