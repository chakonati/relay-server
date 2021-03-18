package handlers

type ActionHandler struct {
	handler *Handler
	request *Request

}

func (a *ActionHandler) Echo(echo string) string {
	return echo
}
