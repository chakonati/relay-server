package actions

import (
	"log"
	"server/defs"
	"server/persistence"

	"github.com/pkg/errors"
)

type KeyExchangeHandler struct {
	setup SetupHandler
}

func (k *KeyExchangeHandler) PublishPreKeyBundle(
	registrationId int, deviceId int, preKeyId int,
	publicPreKey []byte, signedPreKeyId int,
	publicSignedPreKey []byte, signedPreKeySignature []byte,
	identityKey []byte, password string,
) error {
	if !k.setup.IsPasswordValid(password) {
		return InvalidPasswordError()
	}

	log.Println("Received pre-key bundle")
	if err := persistence.KeyExchange.DeleteAllOneTimePreKeys(); err != nil {
		return errors.Wrap(err, "could not delete previous one time pre-keys")
	}

	preKeyBundle := defs.PreKeyBundle{
		RegistrationID:        registrationId,
		DeviceID:              deviceId,
		SignedPreKeyID:        signedPreKeyId,
		PublicSignedPreKey:    publicSignedPreKey,
		SignedPreKeySignature: signedPreKeySignature,
		IdentityKey:           identityKey,
	}
	if err := persistence.KeyExchange.StorePreKeyBundle(preKeyBundle); err != nil {
		return errors.Wrap(err, "failed to store key bundle")
	}

	if err := persistence.KeyExchange.AddOneTimePreKey(defs.OneTimePreKey{
		PreKeyId: preKeyId,
		PreKey:   publicPreKey,
	}); err != nil {
		return errors.Wrap(err, "could not store one time pre-key")
	}
	return nil
}

func (k *KeyExchangeHandler) PreKeyBundle() (
	registrationId int, deviceId int,
	preKeyId *int, preKey *[]byte,
	signedPreKeyId int, publicSignedPreKey []byte,
	signedPreKeySignature []byte,
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
		return
	}

	oneTimePreKey, err := k.OneTimePreKey()
	if err != nil {
		err = errors.Wrap(err, "failed to retrieve one time pre-key")
		return
	}

	registrationId = preKeyBundle.RegistrationID
	deviceId = preKeyBundle.DeviceID
	if oneTimePreKey != nil {
		preKeyId = &oneTimePreKey.PreKeyId
		preKey = &oneTimePreKey.PreKey
	}
	signedPreKeyId = preKeyBundle.SignedPreKeyID
	publicSignedPreKey = preKeyBundle.PublicSignedPreKey
	signedPreKeySignature = preKeyBundle.SignedPreKeySignature
	identityKey = preKeyBundle.IdentityKey
	return
}

func (k *KeyExchangeHandler) PreKeyBundleExists() bool {
	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	return err == nil && keyExists
}

func (k *KeyExchangeHandler) PublishOneTimePreKeys(preKeys []defs.OneTimePreKey, password string) error {
	if !k.setup.IsPasswordValid(password) {
		return InvalidPasswordError()
	}

	for index, preKey := range preKeys {
		if err := persistence.KeyExchange.AddOneTimePreKey(preKey); err != nil {
			return errors.Wrapf(err,
				"could not add one time pre-key at index %d, id %d", index, preKey.PreKeyId)
		}
	}

	return nil
}

func (k *KeyExchangeHandler) OneTimePreKey() (*defs.OneTimePreKey, error) {
	preKey, err := persistence.KeyExchange.NextOneTimePreKey()
	if err != nil {
		return nil, errors.Wrap(err, "could not get a pre key")
	}
	return preKey, nil
}

func (k *KeyExchangeHandler) DeviceId() (int, error) {
	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	if err != nil {
		return 0, errors.Wrap(err, "failed to check if the pre-key bundle already exists")
	}
	if !keyExists {
		return 0, errors.New("pre-key bundle does not exist")
	}

	preKeyBundle, err := persistence.KeyExchange.PreKeyBundle()
	if err != nil {
		err = errors.Wrap(err, "failed to retrieve pre-key bundle")
	}

	return preKeyBundle.DeviceID, nil
}
