package transport

import (
	"avito_internship/internal/auth"
	"avito_internship/internal/repository"
	"net/http"
)

func MapRoutes() {
	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		GetJWT(w, r, func(username, password string) (string, error) {
			return auth.Authenticate(username, password, repository.GetUserIDPassHashOrRegister)
		})
	})
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		GetUserInfo(w, r, repository.GetUserBalanceInventoryLogs)
	})
	http.HandleFunc("/api/sendCoin", func(w http.ResponseWriter, r *http.Request) {
		TransferCoins(w, r, repository.SendCoins)
	})
	http.HandleFunc("/api/buy/", func(w http.ResponseWriter, r *http.Request) {
		BuyItems(w, r, repository.BuyItemsForUser)
	})
}
