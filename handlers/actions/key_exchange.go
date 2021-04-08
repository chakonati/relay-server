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

type KeyExchangePublishPreKeyBundleResponse struct {
	Error error
}

type KeyExchangePublishPreKeyBundleRequest struct {
	RegistrationId        int
	DeviceId              int
	PreKeyId              int
	PublicPreKey          []byte
	SignedPreKeyId        int
	PublicSignedPreKey    []byte
	SignedPreKeySignature []byte
	IdentityKey           []byte
	Password              string
}

func (k *KeyExchangeHandler) PublishPreKeyBundle(
	request KeyExchangePublishPreKeyBundleRequest,
) *KeyExchangePublishPreKeyBundleResponse {
	return &KeyExchangePublishPreKeyBundleResponse{Error: k.publishPreKeyBundle(
		request.RegistrationId, request.DeviceId, request.PreKeyId,
		request.PublicPreKey, request.SignedPreKeyId,
		request.PublicSignedPreKey, request.SignedPreKeySignature,
		request.IdentityKey, request.Password,
	)}
}

func (k *KeyExchangeHandler) publishPreKeyBundle(
	registrationId int, deviceId int, preKeyId int,
	publicPreKey []byte, signedPreKeyId int,
	publicSignedPreKey []byte, signedPreKeySignature []byte,
	identityKey []byte, password string,
) error {
	if !k.setup.isPasswordValid(password).Valid {
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

type KeyExchangePreKeyBundleResponse struct {
	RegistrationId        int
	DeviceId              int
	PreKeyId              *int
	PreKey                *[]byte
	SignedPreKeyId        int
	PublicSignedPreKey    []byte
	SignedPreKeySignature []byte
	IdentityKey           []byte
	Password              string
	Error                 error
}

func (k *KeyExchangeHandler) PreKeyBundle() (response *KeyExchangePreKeyBundleResponse) {
	response = &KeyExchangePreKeyBundleResponse{}

	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	if err != nil {
		response.Error = errors.Wrap(err, "failed to check if the pre-key bundle already exists")
		return
	}
	if !keyExists {
		err = errors.New("pre-key bundle does not exist")
		return
	}

	preKeyBundle, err := persistence.KeyExchange.PreKeyBundle()
	if err != nil {
		response.Error = errors.Wrap(err, "failed to retrieve pre-key bundle")
		return
	}

	oneTimePreKeyResponse := k.OneTimePreKey()
	if oneTimePreKeyResponse.Error != nil {
		response.Error = errors.Wrap(oneTimePreKeyResponse.Error, "failed to retrieve one time pre-key")
		return
	}

	response.RegistrationId = preKeyBundle.RegistrationID
	response.DeviceId = preKeyBundle.DeviceID
	if oneTimePreKeyResponse.OneTimePreKey != nil {
		response.PreKeyId = &oneTimePreKeyResponse.OneTimePreKey.PreKeyId
		response.PreKey = &oneTimePreKeyResponse.OneTimePreKey.PreKey
	}
	response.SignedPreKeyId = preKeyBundle.SignedPreKeyID
	response.PublicSignedPreKey = preKeyBundle.PublicSignedPreKey
	response.SignedPreKeySignature = preKeyBundle.SignedPreKeySignature
	response.IdentityKey = preKeyBundle.IdentityKey
	return
}

type KeyExchangePreKeyBundleExistsResponse struct{ Exists bool }

func (k *KeyExchangeHandler) PreKeyBundleExists() *KeyExchangePreKeyBundleExistsResponse {
	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	return &KeyExchangePreKeyBundleExistsResponse{err == nil && keyExists}
}

type KeyExchangePublishOneTimePreKeysResponse struct{ Error error }

type KeyExchangePublishOneTimePreKeysRequest struct {
	PreKeys  []defs.OneTimePreKey
	Password string
}

func (k *KeyExchangeHandler) PublishOneTimePreKeys(
	request KeyExchangePublishOneTimePreKeysRequest,
) *KeyExchangePublishOneTimePreKeysResponse {
	password := request.Password
	preKeys := request.PreKeys

	if !k.setup.isPasswordValid(password).Valid {
		return &KeyExchangePublishOneTimePreKeysResponse{InvalidPasswordError()}
	}

	for index, preKey := range preKeys {
		if err := persistence.KeyExchange.AddOneTimePreKey(preKey); err != nil {
			return &KeyExchangePublishOneTimePreKeysResponse{errors.Wrapf(err,
				"could not add one time pre-key at index %d, id %d", index, preKey.PreKeyId)}
		}
	}

	return &KeyExchangePublishOneTimePreKeysResponse{}
}

type KeyExchangeOneTimePreKeyResponse struct {
	OneTimePreKey *defs.OneTimePreKey
	Error         error
}

func (k *KeyExchangeHandler) OneTimePreKey() *KeyExchangeOneTimePreKeyResponse {
	preKey, err := persistence.KeyExchange.NextOneTimePreKey()
	if err != nil {
		return &KeyExchangeOneTimePreKeyResponse{nil, errors.Wrap(err, "could not get a pre key")}
	}
	return &KeyExchangeOneTimePreKeyResponse{preKey, nil}
}

type KeyExchangeDeviceIDResponse struct {
	DeviceId int
	Error    error
}

func (k *KeyExchangeHandler) DeviceId() (response *KeyExchangeDeviceIDResponse) {
	response = &KeyExchangeDeviceIDResponse{}

	keyExists, err := persistence.KeyExchange.PreKeyBundleExists()
	if err != nil {
		response.Error = errors.Wrap(err, "failed to check if the pre-key bundle already exists")
	}
	if !keyExists {
		response.Error = errors.New("pre-key bundle does not exist")
	}

	preKeyBundle, err := persistence.KeyExchange.PreKeyBundle()
	if err != nil {
		err = errors.Wrap(err, "failed to retrieve pre-key bundle")
	}

	response.DeviceId = preKeyBundle.DeviceID
	return
}
