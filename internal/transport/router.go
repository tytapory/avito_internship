package transport

import (
	"avito_internship/internal/auth"
	"net/http"
)

func MapRoutes() {
	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		Authentication(w, r, auth.Authenticate)
	})
}
