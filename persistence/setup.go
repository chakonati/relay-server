package persistence

import (
	"log"
	"server/defs"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type SetupDAO struct{}

const setupBucketName = "key_exchange"

var (
	setupPasswordHashKey = []byte("password")
)

var Setup = SetupDAO{}

type PasswordHash struct {
	Hash      []byte
	Algorithm defs.PasswordHashingAlgorithm
}

func (k *SetupDAO) init() error {
	return initBucket(setupBucketName)
}

func (k *SetupDAO) bucket(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket([]byte(setupBucketName))
}

func (k *SetupDAO) StorePasswordHash(hash *PasswordHash) error {
	return setup.db.Update(func(tx *bbolt.Tx) error {
		if err := PutStruct(k.bucket(tx), setupPasswordHashKey, hash); err != nil {
			return errors.Wrap(err, "could not store password hash")
		}
		log.Println("Password hash has been set.")
		return nil
	})
}

func (k *SetupDAO) PasswordHash() (*PasswordHash, error) {
	var passwordHash *PasswordHash
	err := setup.db.View(func(tx *bbolt.Tx) error {
		return GetStruct(k.bucket(tx), setupPasswordHashKey, &passwordHash)
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve password hash")
	}
	return passwordHash, nil
}

func (k *SetupDAO) PasswordExists() (bool, error) {
	bundle, err := k.PasswordHash()
	if err != nil {
		// we return true here just in case the other end doesn't check the error
		// as we don't want this to be overwritten
		log.Println(err)
		return true, errors.Wrap(err, "could not check if password exists")
	}
	return bundle != nil, nil
}
