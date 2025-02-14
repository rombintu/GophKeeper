package storage

import "time"

const (
	Text SecretType = "text"
	Cred SecretType = "cred"
	Data SecretType = "data"
	Card SecretType = "card"
)

type SecretType string
type Payload []byte

type Secret struct {
	Title      string
	SecretType SecretType
	UserID     int64
	Date       time.Time
	Version    int64
	Payload    Payload
}

func NewSecret(title string, secretType SecretType) Secret {
	return Secret{
		Title:      title,
		SecretType: secretType,
		Date:       time.Now(),
	}
}

func (s *Secret) SetPyaload(payload Payload) {
	s.Payload = payload
}

func (s *Secret) SetUserID(userID int64) {
	s.UserID = userID
}

func (s *Secret) IncrementVersion() {
	s.Version++
}
