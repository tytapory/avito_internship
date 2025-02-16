package transport

import (
	"avito_internship/internal/auth"
	"log"
	"net/http"
)

func Run() {
	MapRoutes()
	authMiddleware := Authenticate(http.DefaultServeMux, auth.VerifyJWT)
	log.Fatal(http.ListenAndServe(":8080", authMiddleware))
}
