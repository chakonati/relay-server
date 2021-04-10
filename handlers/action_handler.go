package handlers

import (
	"server/handlers/actions"
	"server/subscriptions"

	"github.com/stoewer/go-strcase"

	"github.com/pkg/errors"
	"gitlab.com/xdevs23/go-reflectutil/reflectutil"
)

const subscriptionSuffix = "Subscription"

type ActionHandler struct {
	handler *Handler
	request *Request

	// Action handlers

	KeyExchange actions.KeyExchangeHandler
	Setup       actions.SetupHandler
	Messaging   actions.MessageHandler
}

type SubscribeRequest struct {
	SubName string
}

type SubscribeResponse struct {
	Error error
}

func (a *ActionHandler) Subscribe(request SubscribeRequest) *SubscribeResponse {
	var subscription subscriptions.Subscription
	err := reflectutil.ExtractByName(
		&subscriptions.Subscriptions,
		strcase.UpperCamelCase(request.SubName)+subscriptionSuffix,
		&subscription,
	)
	if err != nil {
		return &SubscribeResponse{errors.Wrapf(err, "could not find subscription with name %s", request.SubName)}
	}

	subscription.Subscribe(a.handler)

	return &SubscribeResponse{}
}
