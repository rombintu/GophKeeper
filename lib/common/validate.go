package common

import (
	"errors"

	"github.com/rombintu/GophKeeper/internal/proto"
)

// TODO: Email validate
func UserValidate(user *proto.User) error {
	if user.GetEmail() == "" {
		return errors.New("email is required")
	}
	if user.GetHexKeys() == nil {
		return errors.New("fingerprint of keys is required")
	}
	return nil
}
