package persistence

import (
	"log"
	"server/defs"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type KeyExchangeDAO struct{}

const keyExchangeBucketName = "key_exchange"

var (
	keyExchangePreKeyBundleKey = []byte("pre_key_bundle")
)

var KeyExchange = KeyExchangeDAO{}

func (k *KeyExchangeDAO) init() error {
	return setup.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(keyExchangeBucketName))
		if err != nil {
			return errors.Wrap(err, "could not create or get bucket "+keyExchangeBucketName)
		}
		return nil
	})
}

func (k *KeyExchangeDAO) Bucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(keyExchangeBucketName))
}

func (k *KeyExchangeDAO) StorePreKeyBundle(bundle defs.PreKeyBundle) error {
	return setup.db.Update(func(tx *bbolt.Tx) error {
		if err := PutStruct(k.Bucket(tx), keyExchangePreKeyBundleKey, &bundle); err != nil {
			return errors.Wrap(err, "could not store pre-key bundle")
		}
		log.Println("The pre-key bundle is now registered.")
		return nil
	})
}

func (k *KeyExchangeDAO) PreKeyBundle() (*defs.PreKeyBundle, error) {
	var preKeyBundle *defs.PreKeyBundle
	err := setup.db.View(func(tx *bbolt.Tx) error {
		return GetStruct(k.Bucket(tx), keyExchangePreKeyBundleKey, &preKeyBundle)
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve pre-key bundle")
	}
	return preKeyBundle, nil
}

func (k *KeyExchangeDAO) PreKeyBundleExists() (bool, error) {
	bundle, err := k.PreKeyBundle()
	if err != nil {
		// we return true here just in case the other end doesn't check the error
		// as we don't want this to be overwritten
		return true, errors.Wrap(err, "could not check if pre-key bundle exists")
	}
	return bundle != nil, nil
}
