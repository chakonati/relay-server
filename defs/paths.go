package defs

import "path"

const dataDir = "data"

func DataDir() string {
	return dataDir
}

func DatabaseDir() string {
	return path.Join(DataDir(), "databases")
}

// SetupDatabase returns the path to the database containing
// information about this instance and how it is set up.
// Initial settings and long-term settings may be stored here.
func SetupDatabase() string {
	return path.Join(DatabaseDir(), "setup.db")
}

// KeyDatabase returns the path to the database containing
// various types of key material like pre-keys and ephemeral
// keys
func KeyDatabase() string {
	return path.Join(DatabaseDir(), "keys.db")
}

// MessageDatabase returns any kind of encrypted messages
// that are stored temporarily until the recipient retrieves
// them.
func MessageDatabase() string {
	return path.Join(DatabaseDir(), "messages.db")
}
