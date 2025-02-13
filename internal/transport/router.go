package transport

import "net/http"

func MapRoutes() {
	http.HandleFunc("/api/auth", Authentication)
}
