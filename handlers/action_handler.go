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

func (a *ActionHandler) Echo(echo string) string {
	return echo
}

func (a *ActionHandler) Subscribe(subName string) error {
	var subscription subscriptions.Subscription
	err := reflectutil.ExtractByName(
		subscriptions.Subscriptions,
		strcase.UpperCamelCase(subName)+subscriptionSuffix,
		&subscription,
	)
	if err != nil {
		return errors.Wrapf(err, "could not find subscription with name %s", subName)
	}

	subscription.Subscribe(a.handler)

	return nil
}
