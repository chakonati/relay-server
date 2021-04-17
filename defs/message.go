package defs

import "time"

type Message struct {
	ID               uint64    `key:"id"`
	EncryptedMessage []byte    `key:"encryptedMessage"`
	From             string    `key:"from"`
	DeviceID         int       `key:"deviceId"`
	At               time.Time `key:"at"`
}
