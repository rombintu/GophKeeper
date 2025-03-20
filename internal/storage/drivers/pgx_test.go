package drivers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgreSQLContainer(ctx context.Context) (testcontainers.Container, string, error) {
	// Получаем абсолютный путь к директории migrations
	migrationsPath, err := filepath.Abs(filepath.Join("..", "migrations"))
	if err != nil {
		return nil, "", fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	fmt.Println(migrationsPath)

	// Запрос для создания контейнера
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			// Монтирование директории
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: migrationsPath,    // Путь на хосте
				Target: "/tmp/migrations", // Путь в контейнере
			})
		},
	}

	// Запуск контейнера
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	// Получение адреса контейнера
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, "", err
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", err
	}

	// Формирование строки подключения
	connStr := fmt.Sprintf("postgres://test:test@%s:%s/test?sslmode=disable", host, port.Port())

	return postgresContainer, connStr, nil
}

func TestPgxDriver_UserCreate_Docker(t *testing.T) {
	ctx := context.Background()

	// Запуск контейнера PostgreSQL
	postgresContainer, connStr, err := setupPostgreSQLContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(ctx) // Остановка контейнера после теста

	// Инициализация PgxDriver
	d := &PgxDriver{
		dbURL:       connStr,
		serviceName: "test",
	}
	if err := d.Open(ctx); err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	testCtx := context.WithValue(ctx, testKey, true)
	if err := d.Configure(testCtx); err != nil {
		t.Fatal(err)
	}
	// Тестовые данные
	user := &apb.User{Email: "test@example.com", KeyChecksum: []byte("checksum")}

	// Тест
	err = d.UserCreate(ctx, user)
	assert.NoError(t, err, "UserCreate should not return an error")

	// Проверка, что пользователь создан
	foundUser, err := d.UserGet(ctx, user)
	assert.NoError(t, err, "UserGet should not return an error")
	assert.Equal(t, user.Email, foundUser.Email, "Emails should match")
	assert.Equal(t, user.KeyChecksum, foundUser.KeyChecksum, "KeyChecksums should match")
}
func TestPgxDriver_SecretCreateBatch_Docker(t *testing.T) {
	ctx := context.Background()

	// Запуск контейнера PostgreSQL
	postgresContainer, connStr, err := setupPostgreSQLContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(ctx) // Остановка контейнера после теста

	// Инициализация PgxDriver
	d := &PgxDriver{
		dbURL:       connStr,
		serviceName: "test",
	}
	if err := d.Open(ctx); err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}

	testCtx := context.WithValue(ctx, testKey, true)
	if err := d.Configure(testCtx); err != nil {
		t.Fatal(err)
	}
	// Подготовка данных
	user := &apb.User{Email: "test@example.com", KeyChecksum: []byte("checksum")}
	if err := d.UserCreate(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Тестовые данные
	secrets := []*kpb.Secret{
		{
			Title:       "test1",
			SecretType:  kpb.Secret_TEXT,
			UserEmail:   "test@example.com",
			Version:     0,
			HashPayload: "hash1",
			Payload:     []byte("payload1"),
		},
		{
			Title:       "test2",
			SecretType:  kpb.Secret_TEXT,
			UserEmail:   "test@example.com",
			Version:     1,
			HashPayload: "hash2",
			Payload:     []byte("payload2"),
		},
	}

	// Тест
	err = d.SecretCreateBatch(ctx, secrets)
	assert.NoError(t, err, "SecretCreateBatch should not return an error")

	// Проверка, что секреты созданы
	foundSecrets, err := d.SecretList(ctx, "test@example.com")
	assert.NoError(t, err, "SecretList should not return an error")
	assert.Len(t, foundSecrets, 2, "Expected 2 secrets")
}
func TestPgxDriver_Ping(t *testing.T) {
	dbPath := os.Getenv("PGX_DB_PATH")
	ctx := context.Background()
	d := &PgxDriver{
		dbURL:       dbPath,
		serviceName: "test",
	}

	if err := d.Open(ctx); err != nil {
		t.Skipf("Skipping test due to database connection error: %v", err)
		return
	}

	tests := []struct {
		name       string
		monitoring bool
		wantErr    bool
	}{
		{
			name:       "ping_with_monitoring",
			monitoring: true,
			wantErr:    false,
		},
		{
			name:       "ping_without_monitoring",
			monitoring: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Ping(ctx, tt.monitoring); (err != nil) != tt.wantErr {
				t.Errorf("PgxDriver.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
