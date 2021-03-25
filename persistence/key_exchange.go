package persistence

import (
	"log"
	"server/defs"
	"sync"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type KeyExchangeDAO struct{}

const keyExchangeBucketName = "key_exchange"
const keyExchangeOTPKBucketName = "key_exchange_otpk"

var (
	keyExchangePreKeyBundleKey = []byte("pre_key_bundle")
)

var KeyExchange = KeyExchangeDAO{}

func (k *KeyExchangeDAO) init() error {
	if err := initBucket(&setup, keyExchangeBucketName); err != nil {
		return err
	}
	if err := initBucket(&keys, keyExchangeOTPKBucketName); err != nil {
		return err
	}
	return nil
}

func (k *KeyExchangeDAO) bucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(keyExchangeBucketName))
}

func (k *KeyExchangeDAO) otpkBucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(keyExchangeOTPKBucketName))
}

func (k *KeyExchangeDAO) StorePreKeyBundle(bundle defs.PreKeyBundle) error {
	return setup.db.Update(func(tx *bbolt.Tx) error {
		if err := PutStruct(k.bucket(tx), keyExchangePreKeyBundleKey, &bundle); err != nil {
			return errors.Wrap(err, "could not store pre-key bundle")
		}
		log.Println("The pre-key bundle is now registered.")
		return nil
	})
}

func (k *KeyExchangeDAO) PreKeyBundle() (*defs.PreKeyBundle, error) {
	var preKeyBundle *defs.PreKeyBundle
	err := setup.db.View(func(tx *bbolt.Tx) error {
		return GetStruct(k.bucket(tx), keyExchangePreKeyBundleKey, &preKeyBundle)
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

var otpkMut sync.Mutex

func (k *KeyExchangeDAO) AddOneTimePreKey(preKey defs.OneTimePreKey) error {
	return keys.db.Update(func(tx *bbolt.Tx) error {
		if err := k.otpkBucket(tx).Put(IntByteArray(preKey.PreKeyId), preKey.PreKey); err != nil {
			return errors.Wrap(err, "could not add one time pre-key")
		}
		return nil
	})
}

func (k *KeyExchangeDAO) NextOneTimePreKey() (*defs.OneTimePreKey, error) {
	var preKey *defs.OneTimePreKey
	otpkMut.Lock()
	defer otpkMut.Unlock()
	return preKey, keys.db.Update(func(tx *bbolt.Tx) error {
		cursor := k.otpkBucket(tx).Cursor()
		idByt, key := cursor.First()
		if idByt != nil {
			preKey = &defs.OneTimePreKey{
				PreKeyId: ByteArrayInt(idByt),
				PreKey:   key,
			}
			if err := cursor.Delete(); err != nil {
				return errors.Wrap(err, "could not delete retrieved one time pre-key")
			}
		}
		return nil
	})
}
