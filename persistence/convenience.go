package persistence

import (
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
