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
)

func TestPurchaseAndInventory(t *testing.T) {
	baseURL := "http://avito-shop-service-test:8080"
	authURL := baseURL + "/api/auth"
	infoURL := baseURL + "/api/info"
	purchaseURL := baseURL + "/api/buy/t-shirt"

	regPayload := models.AuthRequest{
		Username: "testUserForPurchase",
		Password: "password123",
	}
	jsonPayload, err := json.Marshal(regPayload)
	require.NoError(t, err)

	respReg, err := http.Post(authURL, "application/json", bytes.NewBuffer(jsonPayload))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respReg.StatusCode)
	defer respReg.Body.Close()

	bodyReg, err := io.ReadAll(respReg.Body)
	require.NoError(t, err)

	var regTokenResponse models.AuthResponse
	err = json.Unmarshal(bodyReg, &regTokenResponse)
	require.NoError(t, err)
	require.NotEmpty(t, regTokenResponse.Token)

	token := regTokenResponse.Token

	client := &http.Client{}
	reqPurchase, err := http.NewRequest("GET", purchaseURL, nil)
	require.NoError(t, err)
	reqPurchase.Header.Set("Authorization", "Bearer "+token)

	respPurchase, err := client.Do(reqPurchase)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respPurchase.StatusCode)
	defer respPurchase.Body.Close()

	reqInfo, err := http.NewRequest("GET", infoURL, nil)
	require.NoError(t, err)
	reqInfo.Header.Set("Authorization", "Bearer "+token)

	respInfo, err := client.Do(reqInfo)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respInfo.StatusCode)
	defer respInfo.Body.Close()

	infoBody, err := io.ReadAll(respInfo.Body)
	require.NoError(t, err)

	var infoResponse models.InfoResponse
	err = json.Unmarshal(infoBody, &infoResponse)
	require.NoError(t, err)

	expectedCoins := 1000 - 80
	assert.Equal(t, expectedCoins, infoResponse.Coins)

	var found bool
	for _, item := range infoResponse.Inventory {
		if item.Type == "t-shirt" {
			found = true
			assert.Equal(t, 1, item.Quantity)
		}
	}
	assert.True(t, found)
}
