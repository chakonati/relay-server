package handlers

type ActionHandler struct {
	Handler *Handler
	Request *Request
}

func (a *ActionHandler) Echo(echo string) string {
	return echo
}
