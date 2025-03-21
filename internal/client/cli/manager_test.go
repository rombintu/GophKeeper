package cli

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rombintu/GophKeeper/internal/client/models"
	"github.com/rombintu/GophKeeper/internal/proto/auth"
	"github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage/mocks"
	"github.com/stretchr/testify/assert"
)

func TestManager_Configure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockClientManager(ctrl)
	manager := NewManager()
	manager.SetStore(mockStore)

	// Ожидаем вызов Configure с любым контекстом и возвращаем nil
	mockStore.EXPECT().Configure(gomock.Any()).Return(nil)

	err := manager.Configure()
	assert.NoError(t, err)
}

func TestManager_SecretList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockClientManager(ctrl)
	manager := NewManager()
	manager.SetStore(mockStore)

	// Создаем тестовый профиль
	profile := &Profile{
		user: &auth.User{
			Email: "test@example.com",
		},
	}
	manager.SetProfile(profile)

	// Ожидаем вызов SecretList с конкретным email и возвращаем тестовые данные
	expectedSecrets := []*keeper.Secret{
		{
			Title:       "Test Secret",
			SecretType:  keeper.Secret_TEXT,
			UserEmail:   "test@example.com",
			Payload:     []byte("test payload"),
			HashPayload: "hash",
		},
	}
	mockStore.EXPECT().SecretList(gomock.Any(), "test@example.com").Return(expectedSecrets, nil)

	err := manager.SecretList(context.Background())
	assert.NoError(t, err)
}

func TestManager_SecretCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockClientManager(ctrl)
	manager := NewManager()
	manager.SetStore(mockStore)

	// Создаем тестовый профиль
	profile := &Profile{
		user: &auth.User{
			Email: "test@example.com",
		},
	}
	manager.SetProfile(profile)

	// Создаем тестовый секрет
	secret := models.SecretAdapter(&models.SecretText{
		Text: "Test Secret",
	})

	// Ожидаем вызов SecretCreate с конкретными параметрами и возвращаем nil
	mockStore.EXPECT().SecretCreate(gomock.Any(), gomock.Any()).Return(nil)

	err := manager.SecretCreate(context.Background(), secret)
	assert.NoError(t, err)
}

func TestManager_ConfigSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockClientManager(ctrl)
	manager := NewManager()
	manager.SetStore(mockStore)

	// Ожидаем вызов Set с конкретными параметрами и возвращаем nil
	mockStore.EXPECT().Set(gomock.Any(), []byte("key1"), []byte("value1")).Return(nil)
	mockStore.EXPECT().Set(gomock.Any(), []byte("key2"), []byte("value2")).Return(nil)

	values := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := manager.ConfigSet(context.Background(), values)
	assert.NoError(t, err)
}

func TestManager_ConfigGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockClientManager(ctrl)
	manager := NewManager()
	manager.SetStore(mockStore)

	// Ожидаем вызов Get с конкретным ключом и возвращаем тестовые данные
	mockStore.EXPECT().Get(gomock.Any(), []byte("key1")).Return([]byte("value1"), nil)

	value, err := manager.ConfigGet(context.Background(), "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)
}

// func TestManager_Login(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockStore := mocks.NewMockClientManager(ctrl)
// 	manager := NewManager()
// 	manager.SetStore(mockStore)

// 	// Создаем тестовый профиль
// 	profile := &Profile{
// 		user: &auth.User{
// 			Email:       "test@example.com",
// 			KeyChecksum: []byte("hash"),
// 		},
// 	}
// 	manager.SetProfile(profile)

// 	err := manager.Login(context.Background(), "localhost:3201")
// 	assert.NoError(t, err)
// }

// func TestManager_Register(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockStore := mocks.NewMockClientManager(ctrl)
// 	manager := NewManager()
// 	manager.SetStore(mockStore)

// 	// Создаем тестовый профиль
// 	profile := &Profile{
// 		user: &auth.User{
// 			Email:       "test@example.com",
// 			KeyChecksum: []byte("hash"),
// 		},
// 	}
// 	manager.SetProfile(profile)

// 	err := manager.Register(context.Background(), "localhost:3201")
// 	assert.NoError(t, err)
// }

// func TestManager_Sync(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockStore := mocks.NewMockClientManager(ctrl)
// 	manager := NewManager()
// 	manager.SetStore(mockStore)

// 	// Создаем тестовый профиль
// 	profile := &Profile{
// 		user: &auth.User{
// 			Email: "test@example.com",
// 		},
// 	}
// 	manager.SetProfile(profile)

// 	// Ожидаем вызов SecretGetBatch и возвращаем тестовые данные
// 	expectedSecrets := []*keeper.Secret{
// 		{
// 			Title:       "Test Secret",
// 			SecretType:  keeper.Secret_TEXT,
// 			UserEmail:   "test@example.com",
// 			Payload:     []byte("test payload"),
// 			HashPayload: "hash",
// 		},
// 	}
// 	mockStore.EXPECT().SecretGetBatch(gomock.Any()).Return(expectedSecrets, nil)

// 	// Ожидаем вызов SecretCreateBatch с конкретными параметрами и возвращаем nil
// 	mockStore.EXPECT().SecretCreateBatch(gomock.Any(), expectedSecrets).Return(nil)

// 	// Здесь можно добавить моки для grpc клиента, если нужно протестировать взаимодействие с сервером

// 	err := manager.Sync(context.Background(), "localhost:3203")
// 	assert.NoError(t, err)
// }
