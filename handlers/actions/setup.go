package actions

import (
	"log"
	"math/rand"
	"server/crypto"
	"server/defs"
	"server/persistence"

	"github.com/pkg/errors"
	"go.step.sm/crypto/randutil"
)

type SetupHandler struct{}

const setupPasswordLenMin = 14
const setupPasswordLenMax = 20

func (s *SetupHandler) SetPassword() (string, error) {
	passwordExists, err := persistence.Setup.PasswordExists()
	if err != nil {
		log.Println(err)
		return "", errors.New("could not check if password exists")
	}
	if passwordExists {
		return "", errors.New("password has already been set")
	}

	password, err := randutil.Alphanumeric(setupPasswordLenMin + rand.Intn(setupPasswordLenMax-setupPasswordLenMin))
	if err != nil {
		return "", errors.New("could not generate new password")
	}
	passwordHash := persistence.PasswordHash{
		Algorithm: defs.ScryptBcrypt,
	}
	passwordHash.Hash, err = crypto.ScryptHash([]byte(password))
	if err != nil {
		log.Println(err)
		return "", errors.New("could not create hash for password")
	}
	err = persistence.Setup.StorePasswordHash(&passwordHash)
	if err != nil {
		log.Println(err)
		return "", errors.New("could not store password hash")
	}

	return password, nil
}

func (s *SetupHandler) IsPasswordSetup() bool {
	exists, err := persistence.Setup.PasswordExists()
	if err != nil {
		log.Println(err)
		return true
	}
	return exists
}

func (s *SetupHandler) IsPasswordValid(password string) bool {
	if !s.IsPasswordSetup() {
		return false
	}

	hash, err := persistence.Setup.PasswordHash()
	if err != nil {
		log.Println(err)
		return false
	}

	if hash.Algorithm != defs.ScryptBcrypt {
		log.Println("uh oh")
		return false
	}

	matches, err := crypto.ScryptCompare([]byte(password), hash.Hash)
	if err != nil {
		return false
	}

	return matches
}
