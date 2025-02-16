package auth

import (
	"avito_internship/internal/config"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// Authenticate выполняет вход или регистрирует пользователя.
// Если пользователь найден, проверяет пароль и возвращает JWT если пароль верен.
// Если пользователя нет, регистрирует его и выдает JWT.
func Authenticate(username, password string, GetUserFromDB func(string, string) (int, []byte, error)) (string, error) {
	if username == "" || password == "" || len(username) >= 32 {
		return "", ErrInvalidCredentials
	}
	providedPassHash := getHash(password)
	userID, passHash, err := GetUserFromDB(username, string(providedPassHash))
	if err != nil {
		return "", err
	}
	if isPasswordCorrect([]byte(password), passHash) {
		return getJWT(userID), nil
	} else {
		return "", ErrInvalidCredentials
	}
}

// VerifyJWT проверяет JWT и возвращает userID, если токен валиден.
func VerifyJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.Get().JWTSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidCredentials
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidCredentials
	}
	userID := int(userIDFloat)
	return userID, nil
}

// getHash хэширует пароль с использованием bcrypt.
func getHash(password string) []byte {
	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return passHash
}

// isPasswordCorrect проверяет соответствие пароля и хэша.
func isPasswordCorrect(providedPass, passHashFromDB []byte) bool {
	err := bcrypt.CompareHashAndPassword(passHashFromDB, providedPass)
	if err != nil {
		return false
	}
	return true
}

// getJWT создает JWT-токен для user_id со сроком действия 24 часа.
func getJWT(userID int) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(config.Get().JWTSecret)
	return tokenString
}
