package transport

import (
	"avito_internship/internal/auth"
	"avito_internship/internal/repository"
	"net/http"
)

func MapRoutes() {
	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		GetJWT(w, r, auth.Authenticate)
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
