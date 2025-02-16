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

// TestTransferCoins это сценарий где 3 пользователя регистрируются и начинают обмениваться монетами между собой.
// После каждой транзакции проверка информации /api/info и сравнение с ожидаемыми данными
func TestTransferCoinsAndInfo(t *testing.T) {
	baseURL := "http://avito-shop-service-test:8080"
	authURL := baseURL + "/api/auth"
	infoURL := baseURL + "/api/info"
	transferURL := baseURL + "/api/sendCoin"

	users := []models.AuthRequest{
		{Username: "userA", Password: "passwordA"},
		{Username: "userB", Password: "passwordB"},
		{Username: "userC", Password: "passwordC"},
	}

	tokens := make(map[string]string)
	for _, u := range users {
		jsonPayload, err := json.Marshal(u)
		require.NoError(t, err)

		resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(jsonPayload))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var tokenResp models.AuthResponse
		err = json.Unmarshal(body, &tokenResp)
		require.NoError(t, err)
		require.NotEmpty(t, tokenResp.Token)

		tokens[u.Username] = tokenResp.Token
	}

	for _, u := range users {
		info := getUserInfo(t, infoURL, tokens[u.Username])
		assert.Equal(t, 1000, info.Coins)
		assert.Empty(t, info.CoinHistory.Received)
		assert.Empty(t, info.CoinHistory.Sent)
	}

	// Пользователь A отправляет 200 монет пользователю B
	transferCoins(t, transferURL, tokens["userA"], "userB", 200)

	// После первого перевода:
	// userA: 1000 - 200 = 800
	// userB: 1000 + 200 = 1200
	// userC: 1000
	infoA := getUserInfo(t, infoURL, tokens["userA"])
	infoB := getUserInfo(t, infoURL, tokens["userB"])
	infoC := getUserInfo(t, infoURL, tokens["userC"])

	assert.Equal(t, 800, infoA.Coins, "UserA должен иметь 800 монет после перевода 200")
	assert.Equal(t, 1200, infoB.Coins, "UserB должен иметь 1200 монет после получения 200")
	assert.Equal(t, 1000, infoC.Coins, "UserC должен иметь 1000 монет")

	require.Len(t, infoA.CoinHistory.Sent, 1, "UserA должен иметь 1 отправленную транзакцию")
	assert.Equal(t, "userB", infoA.CoinHistory.Sent[0].User, "Отправленная транзакция у UserA должна быть на userB")
	assert.Equal(t, 200, infoA.CoinHistory.Sent[0].Amount, "Сумма отправленной транзакции у UserA должна быть 200")

	require.Len(t, infoB.CoinHistory.Received, 1, "UserB должен иметь 1 полученную транзакцию")
	assert.Equal(t, "userA", infoB.CoinHistory.Received[0].User, "Полученная транзакция у UserB должна быть от userA")
	assert.Equal(t, 200, infoB.CoinHistory.Received[0].Amount, "Сумма полученной транзакции у UserB должна быть 200")

	// Пользователь B отправляет 300 монет пользователю C
	transferCoins(t, transferURL, tokens["userB"], "userC", 300)

	// После второго перевода:
	// userA: 800
	// userB: 1200 - 300 = 900
	// userC: 1000 + 300 = 1300
	infoA = getUserInfo(t, infoURL, tokens["userA"])
	infoB = getUserInfo(t, infoURL, tokens["userB"])
	infoC = getUserInfo(t, infoURL, tokens["userC"])

	assert.Equal(t, 800, infoA.Coins, "UserA должен иметь 800 монет")
	assert.Equal(t, 900, infoB.Coins, "UserB должен иметь 900 монет после отправки 300")
	assert.Equal(t, 1300, infoC.Coins, "UserC должен иметь 1300 монет после получения 300")

	require.Len(t, infoB.CoinHistory.Sent, 1, "UserB должен иметь 1 отправленную транзакцию")
	assert.Equal(t, "userC", infoB.CoinHistory.Sent[0].User, "Отправленная транзакция у UserB должна быть на userC")
	assert.Equal(t, 300, infoB.CoinHistory.Sent[0].Amount, "Сумма отправленной транзакции у UserB должна быть 300")

	require.Len(t, infoC.CoinHistory.Received, 1, "UserC должен иметь 1 полученную транзакцию")
	assert.Equal(t, "userB", infoC.CoinHistory.Received[0].User, "Полученная транзакция у UserC должна быть от userB")
	assert.Equal(t, 300, infoC.CoinHistory.Received[0].Amount, "Сумма полученной транзакции у UserC должна быть 300")
}

// getUserInfo отправляет GET-запрос на infoURL с указанным токеном и возвращает полную информацию пользователя.
func getUserInfo(t *testing.T, infoURL, token string) models.InfoResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", infoURL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var infoResp models.InfoResponse
	err = json.Unmarshal(body, &infoResp)
	require.NoError(t, err)
	return infoResp
}

// transferCoins отправляет POST-запрос на transferURL для перевода монет.
func transferCoins(t *testing.T, transferURL, token, toUser string, amount int) {
	client := &http.Client{}
	transferData := models.SendCoinRequest{
		ToUser: toUser,
		Amount: amount,
	}
	jsonPayload, err := json.Marshal(transferData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", transferURL, bytes.NewBuffer(jsonPayload))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}
