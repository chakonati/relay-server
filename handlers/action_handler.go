package handlers

import "server/handlers/actions"

type ActionHandler struct {
	handler *Handler
	request *Request

	// Action handlers

	KeyExchange actions.KeyExchangeHandler
	Setup       actions.SetupHandler
}

func (a *ActionHandler) Echo(echo string) string {
	return echo
}
