// Package persistence provides an interface to access a database
package persistence

import (
	"os"
	"server/defs"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

var setup DB
var keys DB
var messages DB

type DB struct {
	db *bbolt.DB
}

func createDirectories() error {
	if err := os.MkdirAll(defs.DatabaseDir(), 0700); err != nil {
		return errors.Wrap(err, "could not create database directory")
	}

	return nil
}

var commonDBOptions *bbolt.Options = nil

func initSetupDB() error {
	setupDB, err := bbolt.Open(defs.SetupDatabase(), 0600, commonDBOptions)
	if err != nil {
		return errors.Wrap(err, "could not open setup database")
	}
	setup = DB{
		db: setupDB,
	}
	return nil
}

func initKeyDB() error {
	keyDB, err := bbolt.Open(defs.KeyDatabase(), 0600, commonDBOptions)
	if err != nil {
		return errors.Wrap(err, "could not open key database")
	}
	keys = DB{
		db: keyDB,
	}
	return nil
}

func initMessageDB() error {
	messageDB, err := bbolt.Open(defs.MessageDatabase(), 0600, commonDBOptions)
	if err != nil {
		return errors.Wrap(err, "could not open message database")
	}
	messages = DB{
		db: messageDB,
	}
	return nil
}

func initDAOs() error {
	if err := Setup.init(); err != nil {
		return errors.Wrap(err, "could not initialize setup persistence")
	}
	if err := KeyExchange.init(); err != nil {
		return errors.Wrap(err, "could not initialize key exchange persistence")
	}

	return nil
}

func InitDatabases() error {
	if err := createDirectories(); err != nil {
		return errors.Wrap(err, "could not initialize databases")
	}
	if err := initSetupDB(); err != nil {
		return errors.Wrap(err, "could not initialize setup DB")
	}
	if err := initKeyDB(); err != nil {
		return errors.Wrap(err, "couold not initialize key DB")
	}
	if err := initMessageDB(); err != nil {
		return errors.Wrap(err, "could not initialize message DB")
	}
	if err := initDAOs(); err != nil {
		return errors.Wrap(err, "could not initialize persistence handlers")
	}
	return nil
}
