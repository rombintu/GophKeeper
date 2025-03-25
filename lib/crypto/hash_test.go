package crypto_test

import (
	"crypto/sha1"
	"encoding/base64"
	"testing"

	"github.com/rombintu/GophKeeper/lib/crypto"
	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Empty input",
			input:    []byte(""),
			expected: getExpectedHash([]byte("")),
		},
		{
			name:     "Simple string",
			input:    []byte("hello world"),
			expected: getExpectedHash([]byte("hello world")),
		},
		{
			name:     "Special characters",
			input:    []byte("!@#$%^&*()"),
			expected: getExpectedHash([]byte("!@#$%^&*()")),
		},
		{
			name:     "Unicode characters",
			input:    []byte("こんにちは世界"),
			expected: getExpectedHash([]byte("こんにちは世界")),
		},
		{
			name:     "Long input",
			input:    []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
			expected: getExpectedHash([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crypto.GetHash(tt.input)
			assert.Equal(t, tt.expected, result, "Hash should match expected value")
		})
	}
}

// getExpectedHash вычисляет ожидаемый хеш для тестов
func getExpectedHash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
