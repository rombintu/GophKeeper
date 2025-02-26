GophKeeper

### Generate secret for services
```bash
go run cmd/auth/main.go --secret
```

### Client configuration
```bash
gpg --full-generate-key # Генерация ключей
gpg --export -a "Имя профиля" > profiles/public-key.asc # Импорт публичного ключа
gpg --export-secret-keys -a "Имя профиля" > profiles/private-key.asc # Импорт приватного ключа
```