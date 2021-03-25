package persistence

import (
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
)

func PutStruct(b *bbolt.Bucket, key []byte, s interface{}) error {
	byt, err := msgpack.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "could not marshal struct to put in database")
	}
	if err = b.Put(key, byt); err != nil {
		return errors.Wrap(err, "could not put marshaled struct into bucket")
	}
	return nil
}

func GetStruct(b *bbolt.Bucket, key []byte, out interface{}) error {
	byt := b.Get(key)
	if byt == nil {
		return nil
	}
	if err := msgpack.Unmarshal(byt, out); err != nil {
		return errors.Wrap(err, "could not unmarshal data from database into struct")
	}
	return nil
}

func initBucket(db *DB, bucketName string) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return errors.Wrap(err, "could not create or get bucket "+bucketName)
		}
		return nil
	})
}

func IntByteArray(i int) []byte {
	byt := make([]byte, 4)
	binary.BigEndian.PutUint32(byt, uint32(i))
	return byt
}

func ByteArrayInt(byt []byte) int {
	return int(binary.BigEndian.Uint32(byt))
}
