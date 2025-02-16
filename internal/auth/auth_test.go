package auth

import (
	"avito_internship/internal/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --------------------
// Инициализация тестов
// --------------------
// TestMain устанавливает заведомо известный секрет для проверки правильности работы авторизации
func TestMain(m *testing.M) {
	config.Get().SetJWTSecret([]byte("secret"))
	code := m.Run()
	os.Exit(code)
}

// Валидная пара пароль - хэш для тестов
var validPasswordHashPair = passHashPair{
	password: []byte("test"),
	hash:     []byte("$2a$04$AtKpu2059mhHcr7I524/pO8/A9ucbSaQSOQFsmkdUsEeS1WR/wYG2"),
}

// Невалидная пара пароль - хэш для тестов
var invalidPasswordHashPair = passHashPair{
	password: []byte("test"),
	hash:     []byte("invalid"),
}

type passHashPair struct {
	password []byte
	hash     []byte
}

// -----------------------
// Тесты isPasswordCorrect
// -----------------------
func TestPasswordValidationValid(t *testing.T) {
	assert.True(t, isPasswordCorrect(validPasswordHashPair.password, validPasswordHashPair.hash))
}

func TestPasswordValidationInvalid(t *testing.T) {
	assert.False(t, isPasswordCorrect(invalidPasswordHashPair.password, invalidPasswordHashPair.hash))
}

func TestPasswordValidationInvalidData(t *testing.T) {
	assert.False(t, isPasswordCorrect(nil, []byte("  ")))
}

// -------------
// Тесты getHash
// -------------
func TestGetHashValid(t *testing.T) {
	assert.True(t, isPasswordCorrect(validPasswordHashPair.password, getHash(string(validPasswordHashPair.password))))
}

func TestGetHashEmptyPass(t *testing.T) {
	assert.True(t, isPasswordCorrect([]byte(""), getHash("")))
}

func TestGetHashNilPass(t *testing.T) {
	assert.True(t, isPasswordCorrect(nil, getHash("")))
}

// ---------------
// Тесты VerifyJWT
// ---------------
func TestVerifyJWTValid(t *testing.T) {
	userID, err := VerifyJWT(validToken.token)
	assert.Equal(t, userID, validToken.expectedUserID)
	assert.NoError(t, err)
}

func TestVerifyJWTInvalid(t *testing.T) {
	userID, err := VerifyJWT(invalidToken.token)
	assert.Equal(t, userID, 0)
	assert.Error(t, err)
}

func TestVerifyJWTExpiredToken(t *testing.T) {
	userID, err := VerifyJWT(expiredToken.token)
	assert.Equal(t, userID, 0)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestVerifyJWTValidButSignedWithOtherKey(t *testing.T) {
	userID, err := VerifyJWT(validOtherKeyToken.token)
	assert.Equal(t, userID, 0)
	assert.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
}

// Валидная пара токен - юзер
var validToken = tokenData{
	token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NTAxMDEyNzEsInVzZXJfaWQiOjF9.CRnIINN74Bqz24D2h1zvo4MO6KZQqh89DemFycCjp6I",
	expectedUserID: 1,
}

// Битый токен
var invalidToken = tokenData{
	token:          "invalid",
	expectedUserID: 1,
}

// Валидный токен, но подписанный другим ключем
var validOtherKeyToken = tokenData{
	token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkyMjMzNzIwMzY4NTQ3NzU4MDcsInVzZXJfaWQiOjF9.nG-6NkMf3v_Z2p-x5lJTfXSFnflUDovuadnzFozbNr8",
	expectedUserID: 1,
}

// Просроченный токен, но с правильным юзером
var expiredToken = tokenData{
	token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3MDEzMzAsInVzZXJfaWQiOjF9.EGzRG9cXZ_3sJkirbPWu2che_fjGAAZNTizhNh0QNXo",
	expectedUserID: 1,
}

type tokenData struct {
	token          string
	expectedUserID int
}

// ------------
// Тесты getJWT
// ------------
func TestGetJWTValidUserID(t *testing.T) {
	token := getJWT(1)
	parsedID, _ := VerifyJWT(token)
	assert.Equal(t, parsedID, 1)
}

func TestGetJWTValidNegativeUserID(t *testing.T) {
	token := getJWT(-1)
	parsedID, _ := VerifyJWT(token)
	assert.Equal(t, parsedID, -1)
}

// --------------------
// Тесты Authentication
// --------------------
func TestAuthenticateValid(t *testing.T) {
	token, err := Authenticate("test", string(validPasswordHashPair.password), validGetUserIDPassHashFromDB)
	assert.NoError(t, err, "Ожидалось что аутентификация пройдет успешно")
	assert.NotEmpty(t, token, "Ожидался валидный токен, так как данные верны")
}

func TestAuthenticateInvalidPassword(t *testing.T) {
	token, err := Authenticate("test", string(invalidPasswordHashPair.password), invalidGetUserIDPassHashFromDB)
	assert.Error(t, err, "Ожидалась ошибка аутентификации из-за неверного пароля")
	assert.Empty(t, token, "Ожидался пустой токен, так как пароль неверный")
}

func TestAuthenticateEmptyUsername(t *testing.T) {
	token, err := Authenticate("", string(validPasswordHashPair.password), validGetUserIDPassHashFromDB)
	assert.Error(t, err, "Ожидалась ошибка аутентификации из-за невалидного логина")
	assert.Empty(t, token, "Ожидался пустой токен, так как аутентификация не пройдена")
}

func TestAuthenticateLongUsername(t *testing.T) {
	token, err := Authenticate("1234567890123456789012345678901234567890", string(validPasswordHashPair.password), validGetUserIDPassHashFromDB)
	assert.Error(t, err, "Ожидалась ошибка аутентификации из-за невалидного логина")
	assert.Empty(t, token, "Ожидался пустой токен, так как аутентификация не пройдена")
}

func TestAuthenticateEmptyPassword(t *testing.T) {
	token, err := Authenticate("test", "", validGetUserIDPassHashFromDB)
	assert.Error(t, err, "Ожидалась ошибка аутентификации из-за невалидного пароля")
	assert.Empty(t, token, "Ожидался пустой токен, так как аутентификация не пройдена")
}

func TestAuthenticateErrorFromRepository(t *testing.T) {
	token, err := Authenticate("test", string(validPasswordHashPair.password), errorGetUserIDPassHashFromDB)
	assert.ErrorIs(t, err, databaseError, "Ожидалась ошибка аутентификации из-за ошибки базы данных")
	assert.Empty(t, token, "Ожидался пустой токен, так как аутентификация не пройдена")
}

// errorGetUserIDPassHashFromDB мок функция, которая возвращает помимо верного хэша еще и ошибку
func errorGetUserIDPassHashFromDB(username string, passwordHash string) (int, []byte, error) {
	return 1, validPasswordHashPair.hash, databaseError
}

// Ошибка, которую вернет errorGetUserIDPassHashFromDB
var databaseError = fmt.Errorf("some error")

// validGetUserIDPassHashFromDB мок функция, которая возвращает заведомо верный хеш пароля для test
func validGetUserIDPassHashFromDB(username string, passwordHash string) (int, []byte, error) {
	return 1, validPasswordHashPair.hash, nil
}

// invalidGetUserIDPassHashFromDB мок функция, которая возвращает заведомо неверный хеш пароля для test
func invalidGetUserIDPassHashFromDB(username string, passwordHash string) (int, []byte, error) {
	return 1, invalidPasswordHashPair.hash, nil
}
