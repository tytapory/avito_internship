package transport

type InfoResponse struct {
	Coins     int `json:"coins"`
	Inventory []struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	} `json:"inventory"`
	CoinHistory struct {
		Received []struct {
			FromUser string `json:"fromUser"`
			Amount   int    `json:"amount"`
		} `json:"received"`
		Sent []struct {
			ToUser string `json:"toUser"`
			Amount int    `json:"amount"`
		} `json:"sent"`
	} `json:"coinHistory"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
