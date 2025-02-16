package e2e

import (
	"avito_internship/internal/models"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRegistration(t *testing.T) {
	apiURL := "http://avito-shop-service-test:8080/api/auth"
	payload := models.AuthRequest{
		Username: "testUserID1",
		Password: "testUserID1",
	}
	jsonPayload, err := json.Marshal(payload)
	require.NoError(t, err)

	respReg, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respReg.StatusCode)
	require.NotNil(t, respReg.Body)
	defer respReg.Body.Close()
	bodyReg, err := io.ReadAll(respReg.Body)
	require.NoError(t, err)
	var regTokenResponse models.AuthResponse
	err = json.Unmarshal(bodyReg, &regTokenResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, regTokenResponse.Token)

	time.Sleep(1 * time.Second) // генерация токена зависит от времени в секундах, ждем, чтобы сгенерить новый токен
	respLogin, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respLogin.StatusCode)
	require.NotNil(t, respLogin.Body)
	defer respLogin.Body.Close()
	loginBody, err := io.ReadAll(respLogin.Body)
	require.NoError(t, err)
	var loginTokenResponse models.AuthResponse
	err = json.Unmarshal(loginBody, &loginTokenResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, loginTokenResponse.Token)
	assert.NotEqual(t, regTokenResponse.Token, loginTokenResponse.Token)

	payloadWrongPass := models.AuthRequest{
		Username: "testUserID1",
		Password: "wrongPass",
	}
	jsonPayloadWrongPass, err := json.Marshal(payloadWrongPass)
	require.NoError(t, err)

	respLoginWrong, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayloadWrongPass))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, respLoginWrong.StatusCode)
	require.NotNil(t, respLoginWrong.Body)
	defer respLoginWrong.Body.Close()
	wrongLoginBody, err := io.ReadAll(respLoginWrong.Body)
	require.NoError(t, err)
	var wrongLoginTokenResponse models.AuthResponse
	err = json.Unmarshal(wrongLoginBody, &wrongLoginTokenResponse)
	require.NoError(t, err)
	assert.Empty(t, wrongLoginTokenResponse.Token)
}
