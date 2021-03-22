package defs

type PasswordHashingAlgorithm int

const (
	Bcrypt PasswordHashingAlgorithm = iota
	ScryptBcrypt
)
