package defs

import "time"

type Message struct {
	ID               uint64    `key:"id"`
	EncryptedMessage []byte    `key:"encryptedMessage"`
	At               time.Time `key:"at"`
}
