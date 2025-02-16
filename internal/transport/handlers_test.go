package transport

import (
	"avito_internship/internal/auth"
	"avito_internship/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ------------------
// Тесты Authenticate
// ------------------
func TestAuthenticateValidToken(t *testing.T) {
	mockVerifyJWT := func(token string) (int, error) {
		return 1, nil
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		assert.Equal(t, 1, userID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer validToken")
	rr := httptest.NewRecorder()

	Authenticate(handler, mockVerifyJWT).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthenticateMissingAuthHeader(t *testing.T) {
	mockVerifyJWT := func(token string) (int, error) {
		return 0, errors.New("empty auth header")
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("GET", "/api/protected", nil)
	rr := httptest.NewRecorder()

	Authenticate(handler, mockVerifyJWT).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAuthenticateInvalidToken(t *testing.T) {
	mockVerifyJWT := func(token string) (int, error) {
		return 0, auth.ErrInvalidCredentials
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidToken")
	rr := httptest.NewRecorder()

	Authenticate(handler, mockVerifyJWT).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// --------------
// Тесты BuyItems
// --------------
func TestBuyItemsSuccess(t *testing.T) {
	mockBuyFunc := func(userID int, item string, quantity int) error {
		return nil
	}

	req := httptest.NewRequest("GET", "/api/buy/t_shirt", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	rr := httptest.NewRecorder()

	BuyItems(rr, req, mockBuyFunc)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestBuyItemsInvalidMethod(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/buy/t_shirt", nil)
	rr := httptest.NewRecorder()

	BuyItems(rr, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestBuyItemsInvalidItem(t *testing.T) {
	mockBuyFunc := func(userID int, item string, quantity int) error {
		return errors.New("purchase failed")
	}

	req := httptest.NewRequest("GET", "/api/buy/invalid_item", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	rr := httptest.NewRecorder()

	BuyItems(rr, req, mockBuyFunc)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// -------------------
// Тесты TransferCoins
// -------------------
func TestTransferCoinsSuccess(t *testing.T) {
	mockTransferFunc := func(fromID, amount int, toUser string) error {
		return nil
	}

	reqBody := `{"toUser": "user1", "amount": 50}`
	req := httptest.NewRequest("POST", "/api/sendCoin", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	rr := httptest.NewRecorder()

	TransferCoins(rr, req, mockTransferFunc)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTransferCoinsInvalidMethod(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sendCoin", nil)
	rr := httptest.NewRecorder()

	TransferCoins(rr, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestTransferCoinsFailedTransfer(t *testing.T) {
	mockTransferFunc := func(fromID, amount int, toUser string) error {
		return errors.New("transfer failed")
	}

	reqBody := `{"toUser": "user2", "amount": 50}`
	req := httptest.NewRequest("POST", "/api/sendCoin", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	rr := httptest.NewRecorder()

	TransferCoins(rr, req, mockTransferFunc)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ------------
// Тесты GetJWT
// ------------
func TestGetJWTSuccess(t *testing.T) {
	mockAuthFunc := func(username, password string) (string, error) {
		return "validToken", nil
	}

	reqBody := `{"username": "validUser", "password": "validPass"}`
	req := httptest.NewRequest("POST", "/api/auth", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	GetJWT(rr, req, mockAuthFunc)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetJWTInvalidCredentials(t *testing.T) {
	mockAuthFunc := func(username, password string) (string, error) {
		return "", auth.ErrInvalidCredentials
	}

	reqBody := `{"username": "invalidUser", "password": "wrongPass"}`
	req := httptest.NewRequest("POST", "/api/auth", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	GetJWT(rr, req, mockAuthFunc)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetJWTInvalidMethod(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/auth", nil)
	rr := httptest.NewRecorder()

	GetJWT(rr, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

// ------------
// Тесты GetUserInfo
// ------------
func TestGetUserInfoSuccess(t *testing.T) {
	mockUserInfoFunc := func(userID int) (models.InfoResponse, error) {
		return models.InfoResponse{Coins: 100}, nil
	}

	req := httptest.NewRequest("GET", "/api/userinfo", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
	rr := httptest.NewRecorder()

	GetUserInfo(rr, req, mockUserInfoFunc)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetUserInfoUserNotFound(t *testing.T) {
	mockUserInfoFunc := func(userID int) (models.InfoResponse, error) {
		return models.InfoResponse{}, errors.New("user not found")
	}

	req := httptest.NewRequest("GET", "/api/userinfo", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", 999))
	rr := httptest.NewRecorder()

	GetUserInfo(rr, req, mockUserInfoFunc)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetUserInfoInvalidMethod(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/userinfo", nil)
	rr := httptest.NewRecorder()

	GetUserInfo(rr, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
