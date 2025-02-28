package models

import (
	"time"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

type Secret struct {
	Title string
	Type  kpb.Secret_SecretType
}

// Text model
type SecretText struct {
	Secret
	Text string
}

func (st *SecretText) Encode() []byte {
	return []byte(st.Text)
}

func (st *SecretText) Title() string {
	return st.Secret.Title
}

func (st *SecretText) Type() kpb.Secret_SecretType {
	return st.Secret.Type
}

// Creds model
type SecretCreds struct {
	Secret
	URL      string
	Login    string
	Password string
}

// Binary model
type SecretBinary struct {
	Secret
	BinaryData []byte
}

// Card model
type SecretCard struct {
	Secret
	Owner      string
	ExpireDate time.Time
	Number     string
	Code       int16
}

type SecretAdapter interface {
	Encode() []byte
	Title() string
	Type() kpb.Secret_SecretType
}
