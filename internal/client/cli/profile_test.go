package cli

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAppDirName = "gophkeeper"
)

func TestNewProfile(t *testing.T) {

	// Удаляем тестовую директорию, если она существует
	homeDirPath, err := getHomeDir()
	require.NoError(t, err)
	testProjectDirPath := path.Join(homeDirPath, testAppDirName)
	defer func() {
		_ = os.RemoveAll(testProjectDirPath)
	}()

	// Создаем новый профиль
	profile := NewProfile()

	// Проверяем, что пути к файлам инициализированы корректно
	assert.Equal(t, path.Join(testProjectDirPath, "config.json"), profile.GetConfFilePath())
	assert.Equal(t, path.Join(testProjectDirPath, "bolt.db"), profile.dbPath)

	// Проверяем, что файлы были созданы
	_, err = os.Stat(profile.GetConfFilePath())
	assert.NoError(t, err)

	_, err = os.Stat(profile.dbPath)
	assert.NoError(t, err)
}

func TestGetDriverPath(t *testing.T) {
	// Создаем новый профиль
	profile := NewProfile()

	// Проверяем, что путь к драйверу формируется корректно
	expectedPath := fmt.Sprintf("bolt://%s", profile.dbPath)
	assert.Equal(t, expectedPath, profile.GetDriverPath())
}

func TestGetHomeDir(t *testing.T) {
	// Получаем домашнюю директорию
	homeDir, err := getHomeDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, homeDir)
}

func TestCheckFileExists(t *testing.T) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "test_file")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	// Проверяем, что файл существует
	checkFileExists(tmpFile.Name())

	// Удаляем файл и проверяем, что он был создан заново
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	checkFileExists(tmpFile.Name())

	_, err = os.Stat(tmpFile.Name())
	assert.NoError(t, err)
}
