package models

import (
	"encoding/json"
	"log/slog"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

type Secret struct {
	Title string
	Type  kpb.Secret_SecretType
}

func (s *Secret) GetTitle() string {
	return s.Title
}

func (s *Secret) GetType() kpb.Secret_SecretType {
	return s.Type
}

// Text model
type SecretText struct {
	Secret
	Text string
}

func (st *SecretText) Payload() []byte {
	return []byte(st.Text)
}

type Creds struct {
	URL      string
	Login    string
	Password string
}

// Creds model
type SecretCreds struct {
	Secret
	Creds Creds
}

func (sc *SecretCreds) Payload() []byte {
	data, err := json.Marshal(sc.Creds)
	if err != nil {
		slog.Warn("failed marshal data", slog.String("error", err.Error()))
		data = []byte{}
	}
	return data
}

// Binary model
type SecretBinary struct {
	Secret
	BinaryData []byte
}

func (st *SecretBinary) Payload() []byte {
	return st.BinaryData
}

type Card struct {
	Owner      string
	ExpireDate string
	Number     string
	Code       string
}

// Card model
type SecretCard struct {
	Secret
	Card Card
}

func (sc *SecretCard) Payload() []byte {
	data, err := json.Marshal(sc.Card)
	if err != nil {
		slog.Warn("failed marshal data", slog.String("error", err.Error()))
		data = []byte{}
	}
	return data
}

type SecretAdapter interface {
	Payload() []byte
	GetTitle() string
	GetType() kpb.Secret_SecretType
}
