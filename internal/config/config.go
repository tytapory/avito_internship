package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"sync"
)

var (
	once sync.Once
	cfg  *Config
)

// Config - структура для хранения конфигурации сервиса
type Config struct {
	ServerPort   string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	DatabaseHost string
	JWTSecret    []byte
}

// Get загружает конфигурацию из переменных окружения (только при первом вызове)
// и возвращает указатель на структуру Config
func Get() *Config {
	once.Do(func() {
		cfg = &Config{
			ServerPort:   getEnv("SERVER_PORT", "8080", os.LookupEnv),
			DatabasePort: getEnv("DATABASE_PORT", "5432", os.LookupEnv),
			DatabaseUser: getEnv("DATABASE_USER", "postgres", os.LookupEnv),
			DatabasePass: getEnv("DATABASE_PASSWORD", "password", os.LookupEnv),
			DatabaseName: getEnv("DATABASE_NAME", "mydb", os.LookupEnv),
			DatabaseHost: getEnv("DATABASE_HOST", "localhost", os.LookupEnv),
			JWTSecret:    []byte(getEnv("JWT_SECRET", generateJWTSecret(), os.LookupEnv)),
		}
	})
	return cfg
}

// getEnv получает значение переменной окружения по ключу.
// Если переменная не задана, возвращает значение по умолчанию.
func getEnv(key, fallback string, getEnvFunc func(string) (string, bool)) string {
	if value, ok := getEnvFunc(key); ok {
		return value
	}
	return fallback
}

// generateJWTSecret генерирует ключ для jwt токенов
func generateJWTSecret() string {
	secret := make([]byte, 32)
	rand.Read(secret)
	return base64.StdEncoding.EncodeToString(secret)
}

// SetJWTSecret позволяет устанавливать известный секрет для тестирования
func (c *Config) SetJWTSecret(secret []byte) {
	c.JWTSecret = secret
}
