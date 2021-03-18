package actions

import "log"

type KeyExchangeHandler struct {
}

func (k *KeyExchangeHandler) PublishIdentityKey(identityKey []byte) error {
	log.Println("Received identity key", identityKey)
	return nil
}
