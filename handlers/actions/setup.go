package actions

import (
	"log"
	"math/rand"
	"server/crypto"
	"server/defs"
	"server/persistence"
	"time"

	"github.com/pkg/errors"
)

type SetupHandler struct{}

func (s *SetupHandler) SetPassword(oldPassword string, password string) error {
	passwordExists, err := persistence.Setup.PasswordExists()
	if err != nil {
		log.Println(err)
		return errors.New("could not check if password exists")
	}
	if passwordExists {
		// first check if the old password matches
		currentHash, err := persistence.Setup.PasswordHash()
		if err != nil {
			log.Println(err)
			return errors.New("could not check current password")
		}
		matches, err := crypto.ScryptCompare([]byte(oldPassword), currentHash.Hash)
		if err != nil {
			log.Println(err)
			return errors.New("could not compare passwords")
		}
		if !matches {
			time.Sleep(2*time.Second + time.Duration(rand.Intn(1000))*time.Millisecond)
			return errors.New("old password and hash does not match")
		}
	}
	passwordHash := persistence.PasswordHash{
		Algorithm: defs.ScryptBcrypt,
	}
	passwordHash.Hash, err = crypto.ScryptHash([]byte(password))
	if err != nil {
		log.Println(err)
		return errors.New("could not create hash for password")
	}
	err = persistence.Setup.StorePasswordHash(&passwordHash)
	if err != nil {
		log.Println(err)
		return errors.New("could not store password hash")
	}

	return nil
}
