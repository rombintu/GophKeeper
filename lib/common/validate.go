package common

import (
	"errors"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
)

// TODO: Email validate
func UserValidate(user *apb.User) error {
	if user.GetEmail() == "" {
		return errors.New("email is required")
	}
	if user.GetKeyChecksum() == nil {
		return errors.New("fingerprint of keys is required")
	}
	return nil
}
