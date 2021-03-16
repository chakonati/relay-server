package session

import (
	"nhooyr.io/websocket"
)

type State int

const (
	Created State = iota
)

type Session struct {
	State State
	Conn  *websocket.Conn
}
