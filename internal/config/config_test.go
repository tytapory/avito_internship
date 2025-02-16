package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// ------------
// Тесты getEnv
// ------------
func TestGetEnvExists(t *testing.T) {
	value := getEnv("SERVER_PORT", "8080", mockGetEnv)
	assert.Equal(t, value, "8888")
}

func TestGetEnvDoesNotExists(t *testing.T) {
	value := getEnv("DATABASE_USER", "test", mockGetEnv)
	assert.Equal(t, value, "test")
}

func TestGetEnvEmpty(t *testing.T) {
	value := getEnv("", "test", mockGetEnv)
	assert.Equal(t, value, "test")
}

// -----------------------
// Тесты generateJWTSecret
// -----------------------
func TestGenerateJWTSecretUniqueness(t *testing.T) {
	secret1 := generateJWTSecret()
	secret2 := generateJWTSecret()
	assert.NotEqual(t, secret1, secret2)
}

// mockGetEnv возвращает корректные значения ключей SERVER_PORT и DATABASE_NAME а для остальных значений
// имитирует ненайденное значение
func mockGetEnv(key string) (string, bool) {
	if key == "SERVER_PORT" {
		return "8888", true
	}
	if key == "DATABASE_NAME" {
		return "test", true
	}
	return "", false
}
