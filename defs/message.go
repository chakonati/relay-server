package defs

import "time"

type Message struct {
	ID               uint64
	EncryptedMessage []byte
	At               time.Time
}
