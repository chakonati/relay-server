package actions

import (
	"log"
	"server/defs"
	"server/persistence"

	"github.com/pkg/errors"
)

type KeyExchangeHandler struct {
}

func (k *KeyExchangeHandler) PublishPreKeyBundle(
	registrationId int, deviceId int, preKeyId int,
	publicPreKey []byte, signedPreKeyId int,
	publicSignedPreKey []byte, signedPreKeySignature []byte,
	identityKey []byte,
) error {
	log.Println("Received pre-key bundle")
	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	if err != nil {
		return errors.Wrap(err, "failed to check if the pre-key bundle already exists")
	}
	if keyExists {
		return errors.New("pre-key bundle already exists, cannot overwrite it")
	}

	preKeyBundle := defs.PreKeyBundle{
		RegistrationID:        registrationId,
		DeviceID:              deviceId,
		PreKeyID:              preKeyId,
		PublicPreKey:          publicPreKey,
		SignedPreKeyID:        signedPreKeyId,
		PublicSignedPreKey:    publicSignedPreKey,
		SignedPreKeySignature: signedPreKeySignature,
		IdentityKey:           identityKey,
	}
	if err := persistence.KeyExchange.StorePreKeyBundle(preKeyBundle); err != nil {
		return errors.Wrap(err, "failed to store key bundle")
	}
	return nil
}

func (k *KeyExchangeHandler) PreKeyBundle() (
	registrationId int, deviceId int, preKeyId int,
	publicPreKey []byte, signedPreKeyId int,
	publicSignedPreKey []byte, signedPreKeySignature []byte,
	identityKey []byte, err error,
) {
	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	if err != nil {
		err = errors.Wrap(err, "failed to check if the pre-key bundle already exists")
		return
	}
	if !keyExists {
		err = errors.New("pre-key bundle does not exist")
		return
	}

	preKeyBundle, err := persistence.KeyExchange.PreKeyBundle()
	if err != nil {
		err = errors.Wrap(err, "failed to retrieve pre-key bundle")
	}

	registrationId = preKeyBundle.RegistrationID
	deviceId = preKeyBundle.DeviceID
	preKeyId = preKeyBundle.PreKeyID
	publicPreKey = preKeyBundle.PublicPreKey
	signedPreKeyId = preKeyBundle.SignedPreKeyID
	publicSignedPreKey = preKeyBundle.PublicSignedPreKey
	signedPreKeySignature = preKeyBundle.SignedPreKeySignature
	identityKey = preKeyBundle.IdentityKey
	return
}
