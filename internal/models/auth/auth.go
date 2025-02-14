package auth

type User struct {
	Email   string
	HexKeys []byte
}

// Создание пользователя с использованием отпечатка его ключей
func NewUser(email string, xKeys []byte) User {
	return User{
		Email:   email,
		HexKeys: xKeys,
	}
}

func (u User) GetHexKeys() []byte {
	return u.HexKeys
}
