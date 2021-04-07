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

type SetupSetPasswordResponse struct {
	Password string
	Error    error
}

func (s *SetupHandler) SetPassword() (response *SetupSetPasswordResponse) {
	response = &SetupSetPasswordResponse{}

	passwordExists, err := persistence.Setup.PasswordExists()
	if err != nil {
		log.Println(err)
		response.Error = errors.New("could not check if password exists")
		return
	}
	if passwordExists {
		response.Error = errors.New("password has already been set")
		return
	}

	password, err := randutil.Alphanumeric(setupPasswordLenMin + rand.Intn(setupPasswordLenMax-setupPasswordLenMin))
	if err != nil {
		response.Error = errors.New("could not generate new password")
		return
	}
	passwordHash := persistence.PasswordHash{
		Algorithm: defs.ScryptBcrypt,
	}
	passwordHash.Hash, err = crypto.ScryptHash([]byte(password))
	if err != nil {
		log.Println(err)
		response.Error = errors.New("could not create hash for password")
		return
	}
	err = persistence.Setup.StorePasswordHash(&passwordHash)
	if err != nil {
		log.Println(err)
		response.Error = errors.New("could not store password hash")
		return
	}

	response.Password = password
	return
}

type MessagingIsPasswordSetupResponse struct{ IsSetUp bool }

func (s *SetupHandler) IsPasswordSetup() *MessagingIsPasswordSetupResponse {
	exists, err := persistence.Setup.PasswordExists()
	if err != nil {
		log.Println(err)
		return &MessagingIsPasswordSetupResponse{false}
	}
	return &MessagingIsPasswordSetupResponse{exists}
}

type MessagingIsPasswordValidResponse struct{ Valid bool }

func (s *SetupHandler) IsPasswordValid(password string) *MessagingIsPasswordValidResponse {
	if !s.IsPasswordSetup().IsSetUp {
		return &MessagingIsPasswordValidResponse{false}
	}

	hash, err := persistence.Setup.PasswordHash()
	if err != nil {
		log.Println(err)
		return &MessagingIsPasswordValidResponse{false}
	}

	if hash.Algorithm != defs.ScryptBcrypt {
		log.Println("uh oh")
		return &MessagingIsPasswordValidResponse{false}
	}

	matches, err := crypto.ScryptCompare([]byte(password), hash.Hash)
	if err != nil {
		return &MessagingIsPasswordValidResponse{false}
	}

	return &MessagingIsPasswordValidResponse{matches}
}
