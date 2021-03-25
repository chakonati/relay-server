package actions

import "github.com/pkg/errors"

func InvalidPasswordError() error {
	return errors.New("specified password can't be used")
}
